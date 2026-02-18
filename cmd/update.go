package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

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

func runUpdate() error {
	latest, err := fetchLatestVersion()
	if err != nil {
		return fmt.Errorf("checking latest version: %w", err)
	}

	current := Version
	if !isNewer(latest, current) {
		fmt.Printf("Already up to date (v%s)\n", current)
		return nil
	}

	fmt.Printf("Updating v%s -> v%s\n", current, latest)

	binPath, err := executablePath()
	if err != nil {
		return fmt.Errorf("locating binary: %w", err)
	}

	url := archiveURL(latest)
	if err := downloadAndReplace(url, binPath); err != nil {
		return err
	}

	fmt.Printf("Updated to v%s\n", latest)
	return nil
}

// fetchLatestVersion queries the GitHub releases API and returns the latest
// version string (without "v" prefix).
func fetchLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/Cloverhound/webex-cli/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return strings.TrimPrefix(release.TagName, "v"), nil
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

// archiveURL builds the download URL for a given version, matching the
// GoReleaser naming template: webex-cli_{VERSION}_{OS}_{ARCH}.tar.gz (.zip on Windows).
func archiveURL(version string) string {
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf(
		"https://github.com/Cloverhound/webex-cli/releases/download/v%s/webex-cli_%s_%s_%s.%s",
		version, version, runtime.GOOS, runtime.GOARCH, ext,
	)
}

// downloadAndReplace downloads the archive from url, extracts the binary, and
// atomically replaces the binary at binaryPath.
func downloadAndReplace(url, binaryPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading download: %w", err)
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
