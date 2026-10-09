package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update webex-cli to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdate()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

// versionPinEnv pins the version that install and update use, skipping lookup.
const versionPinEnv = "WEBEX_CLI_VERSION"

// Release locations, variables so tests can point them at a local server. The
// lookup avoids api.github.com where possible: sandboxes and shared CI hosts
// often block it or exhaust its unauthenticated rate limit, while release
// downloads and the Go module proxy stay reachable.
var (
	releaseBaseURL     = "https://github.com/Cloverhound/webex-cli/releases"
	goProxyLatestURL   = "https://proxy.golang.org/github.com/!cloverhound/webex-cli/@latest"
	githubAPILatestURL = "https://api.github.com/repos/Cloverhound/webex-cli/releases/latest"
)

var lookupClient = &http.Client{Timeout: 30 * time.Second}

func runUpdate() error {
	current := Version
	latest := strings.TrimPrefix(os.Getenv(versionPinEnv), "v")
	if latest != "" {
		if latest == current {
			fmt.Printf("Already at pinned version v%s ($%s)\n", current, versionPinEnv)
			return nil
		}
	} else {
		var err error
		latest, err = fetchLatestVersion()
		if err != nil {
			return fmt.Errorf("checking latest version: %w", err)
		}
		if !isNewer(latest, current) {
			fmt.Printf("Already up to date (v%s)\n", current)
			return nil
		}
	}

	fmt.Printf("Updating v%s -> v%s\n", current, latest)

	binPath, err := executablePath()
	if err != nil {
		return fmt.Errorf("locating binary: %w", err)
	}

	if err := downloadAndReplace(latest, binPath); err != nil {
		return err
	}

	fmt.Printf("Updated to v%s\n", latest)

	// Skills are embedded in the binary, so the new binary must run the check
	// to install the skill that matches it.
	check := exec.Command(binPath, "post-install", "--skills-only")
	check.Stdin, check.Stdout, check.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := check.Run(); err != nil {
		fmt.Printf("Warning: skill update check failed: %v\n", err)
	}

	return nil
}

// fetchLatestVersion returns the latest release version (without "v"), trying
// the release checksums file, then the Go module proxy, then the GitHub API.
func fetchLatestVersion() (string, error) {
	type source struct {
		url    string
		lookup func(string) (string, error)
	}
	sources := []source{
		{releaseBaseURL + "/latest/download/checksums.txt", latestFromChecksums},
		{goProxyLatestURL, latestFromGoProxy},
		{githubAPILatestURL, latestFromGitHubAPI},
	}

	var failures []string
	for _, s := range sources {
		v, err := s.lookup(s.url)
		if err == nil && v != "" {
			return v, nil
		}
		if err == nil {
			err = fmt.Errorf("no version found")
		}
		failures = append(failures, fmt.Sprintf("  %s: %v", s.url, err))
	}
	return "", fmt.Errorf("could not determine the latest version; tried:\n%s\nPin a version instead: %s=x.y.z webex update",
		strings.Join(failures, "\n"), versionPinEnv)
}

func latestFromChecksums(url string) (string, error) {
	body, err := httpGetBytes(url, nil)
	if err != nil {
		return "", err
	}
	return versionFromChecksums(body, runtime.GOOS, runtime.GOARCH)
}

func latestFromGoProxy(url string) (string, error) {
	body, err := httpGetBytes(url, nil)
	if err != nil {
		return "", err
	}
	var info struct {
		Version string `json:"Version"`
	}
	if err := json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	return strings.TrimPrefix(info.Version, "v"), nil
}

func latestFromGitHubAPI(url string) (string, error) {
	header := http.Header{}
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		header.Set("Authorization", "Bearer "+tok)
	}
	body, err := httpGetBytes(url, header)
	if err != nil {
		return "", err
	}
	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}
	return strings.TrimPrefix(release.TagName, "v"), nil
}

func httpGetBytes(url string, header http.Header) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := lookupClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// archiveName matches the GoReleaser name template:
// webex-cli_{VERSION}_{OS}_{ARCH}.tar.gz (.zip on Windows).
func archiveName(version, goos, goarch string) string {
	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf("webex-cli_%s_%s_%s.%s", version, goos, goarch, ext)
}

// versionFromChecksums finds the archive for goos/goarch in a checksums.txt
// and returns the version embedded in its file name.
func versionFromChecksums(data []byte, goos, goarch string) (string, error) {
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	suffix := "_" + goos + "_" + goarch + ext
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		name := strings.TrimPrefix(fields[1], "*")
		if strings.HasPrefix(name, "webex-cli_") && strings.HasSuffix(name, suffix) {
			return strings.TrimSuffix(strings.TrimPrefix(name, "webex-cli_"), suffix), nil
		}
	}
	return "", fmt.Errorf("no %s archive listed", strings.TrimPrefix(suffix, "_"))
}

// checksumFor returns the lowercase hex SHA-256 listed for name.
func checksumFor(data []byte, name string) (string, error) {
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == name {
			return strings.ToLower(fields[0]), nil
		}
	}
	return "", fmt.Errorf("%s is not listed in checksums.txt", name)
}

func verifyChecksum(data []byte, want string) error {
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if got != strings.ToLower(want) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", want, got)
	}
	return nil
}

// isNewer returns true if latest is a higher semver than current.
func isNewer(latest, current string) bool {
	parse := func(v string) [3]int {
		var parts [3]int
		for i, s := range strings.SplitN(v, ".", 3) {
			parts[i], _ = strconv.Atoi(s)
		}
		return parts
	}
	l, c := parse(latest), parse(current)
	for i := range 3 {
		if l[i] > c[i] {
			return true
		}
		if l[i] < c[i] {
			return false
		}
	}
	return false
}

// executablePath returns the resolved path to the running binary.
func executablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exe)
}

// downloadRelease downloads the archive for version and this platform and
// verifies it against the release's checksums.txt.
func downloadRelease(version string) ([]byte, error) {
	base := fmt.Sprintf("%s/download/v%s/", releaseBaseURL, version)
	sums, err := httpGetBytes(base+"checksums.txt", nil)
	if err != nil {
		return nil, fmt.Errorf("downloading checksums.txt: %w", err)
	}
	name := archiveName(version, runtime.GOOS, runtime.GOARCH)
	want, err := checksumFor(sums, name)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(base + name)
	if err != nil {
		return nil, fmt.Errorf("downloading release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading download: %w", err)
	}
	if err := verifyChecksum(body, want); err != nil {
		return nil, fmt.Errorf("%s: %w; refusing to install", name, err)
	}
	return body, nil
}

// downloadAndReplace downloads and verifies the release archive, extracts the
// binary, and atomically replaces the binary at binaryPath.
func downloadAndReplace(version, binaryPath string) error {
	body, err := downloadRelease(version)
	if err != nil {
		return err
	}

	var bin []byte
	if runtime.GOOS == "windows" {
		bin, err = extractFromZip(body)
	} else {
		bin, err = extractFromTarGz(body)
	}
	if err != nil {
		return fmt.Errorf("extracting binary: %w", err)
	}

	// Write to a temp file in the same directory so os.Rename is atomic.
	dir := filepath.Dir(binaryPath)
	tmp, err := os.CreateTemp(dir, "webex-update-*")
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied writing to %s — try: sudo webex update", dir)
		}
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(bin); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Chmod(0755); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("setting permissions: %w", err)
	}
	tmp.Close()

	// On Windows, rename the old binary out of the way first.
	if runtime.GOOS == "windows" {
		oldPath := binaryPath + ".old"
		os.Remove(oldPath)
		if err := os.Rename(binaryPath, oldPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("renaming old binary: %w", err)
		}
	}

	if err := os.Rename(tmpPath, binaryPath); err != nil {
		os.Remove(tmpPath)
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied replacing %s — try: sudo webex update", binaryPath)
		}
		return fmt.Errorf("replacing binary: %w", err)
	}

	return nil
}

// extractFromTarGz extracts the "webex" binary from a .tar.gz archive.
func extractFromTarGz(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if filepath.Base(hdr.Name) == "webex" {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("binary not found in archive")
}

// extractFromZip extracts the "webex.exe" binary from a .zip archive.
func extractFromZip(data []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, f := range zr.File {
		if filepath.Base(f.Name) == "webex.exe" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("binary not found in archive")
}
