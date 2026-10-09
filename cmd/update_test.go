package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
)

const sampleChecksums = `0123abcd  webex-cli_0.16.0_darwin_amd64.tar.gz
4567ef01  webex-cli_0.16.0_darwin_arm64.tar.gz
89ab2345  webex-cli_0.16.0_linux_amd64.tar.gz
cdef6789  webex-cli_0.16.0_linux_arm64.tar.gz
0a1b2c3d  webex-cli_0.16.0_windows_amd64.zip
`

func TestVersionFromChecksums(t *testing.T) {
	for _, tc := range []struct{ goos, goarch, want string }{
		{"linux", "amd64", "0.16.0"},
		{"darwin", "arm64", "0.16.0"},
		{"windows", "amd64", "0.16.0"},
	} {
		got, err := versionFromChecksums([]byte(sampleChecksums), tc.goos, tc.goarch)
		if err != nil || got != tc.want {
			t.Errorf("%s/%s: got %q, %v; want %q", tc.goos, tc.goarch, got, err, tc.want)
		}
	}
	if _, err := versionFromChecksums([]byte(sampleChecksums), "windows", "arm64"); err == nil {
		t.Error("windows/arm64: expected an error for a missing archive")
	}
}

func TestChecksumFor(t *testing.T) {
	got, err := checksumFor([]byte(sampleChecksums), "webex-cli_0.16.0_linux_amd64.tar.gz")
	if err != nil || got != "89ab2345" {
		t.Errorf("got %q, %v", got, err)
	}
	if _, err := checksumFor([]byte(sampleChecksums), "webex-cli_0.16.0_linux_amd64.tar"); err == nil {
		t.Error("expected an error for an unlisted name")
	}
}

func TestVerifyChecksum(t *testing.T) {
	data := []byte("archive bytes")
	sum := sha256.Sum256(data)
	good := hex.EncodeToString(sum[:])
	if err := verifyChecksum(data, strings.ToUpper(good)); err != nil {
		t.Errorf("valid checksum rejected: %v", err)
	}
	if err := verifyChecksum([]byte("tampered"), good); err == nil {
		t.Error("tampered data accepted")
	}
}

// setReleaseURLs points the release lookups at srv for the duration of a test.
func setReleaseURLs(t *testing.T, srv *httptest.Server) {
	t.Helper()
	oldRelease, oldProxy, oldAPI := releaseBaseURL, goProxyLatestURL, githubAPILatestURL
	releaseBaseURL = srv.URL + "/releases"
	goProxyLatestURL = srv.URL + "/proxy/@latest"
	githubAPILatestURL = srv.URL + "/api/latest"
	t.Cleanup(func() {
		releaseBaseURL, goProxyLatestURL, githubAPILatestURL = oldRelease, oldProxy, oldAPI
	})
}

func TestFetchLatestVersionFallbackChain(t *testing.T) {
	archive := archiveName("0.17.0", runtime.GOOS, runtime.GOARCH)
	tests := []struct {
		name     string
		handlers map[string]http.HandlerFunc
		want     string
	}{
		{
			name: "checksums",
			handlers: map[string]http.HandlerFunc{
				"/releases/latest/download/checksums.txt": func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "abc  %s\n", archive)
				},
			},
			want: "0.17.0",
		},
		{
			name: "go proxy",
			handlers: map[string]http.HandlerFunc{
				"/proxy/@latest": func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprint(w, `{"Version":"v0.17.1"}`)
				},
			},
			want: "0.17.1",
		},
		{
			name: "github api with token",
			handlers: map[string]http.HandlerFunc{
				"/api/latest": func(w http.ResponseWriter, r *http.Request) {
					if r.Header.Get("Authorization") != "Bearer gh-test" {
						http.Error(w, "forbidden", http.StatusForbidden)
						return
					}
					fmt.Fprint(w, `{"tag_name":"v0.17.2"}`)
				},
			},
			want: "0.17.2",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			for p, h := range tc.handlers {
				mux.HandleFunc(p, h)
			}
			srv := httptest.NewServer(mux)
			defer srv.Close()
			setReleaseURLs(t, srv)
			t.Setenv("GITHUB_TOKEN", "gh-test")

			got, err := fetchLatestVersion()
			if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
		})
	}
}

func TestFetchLatestVersionReportsEveryURL(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	defer srv.Close()
	setReleaseURLs(t, srv)

	_, err := fetchLatestVersion()
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"checksums.txt", "/proxy/@latest", "/api/latest", versionPinEnv + "="} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error does not mention %q:\n%v", want, err)
		}
	}
}

func TestDownloadReleaseVerifiesChecksum(t *testing.T) {
	archive := archiveName("0.17.0", runtime.GOOS, runtime.GOARCH)
	payload := []byte("real archive")
	sum := sha256.Sum256(payload)

	for _, tc := range []struct {
		name    string
		served  []byte
		listed  bool
		wantErr string
	}{
		{"match", payload, true, ""},
		{"tampered", []byte("evil archive"), true, "checksum mismatch"},
		{"unlisted", payload, false, "not listed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/releases/download/v0.17.0/checksums.txt", func(w http.ResponseWriter, r *http.Request) {
				if tc.listed {
					fmt.Fprintf(w, "%s  %s\n", hex.EncodeToString(sum[:]), archive)
				}
			})
			mux.HandleFunc("/releases/download/v0.17.0/"+archive, func(w http.ResponseWriter, r *http.Request) {
				w.Write(tc.served)
			})
			srv := httptest.NewServer(mux)
			defer srv.Close()
			setReleaseURLs(t, srv)

			got, err := downloadRelease("0.17.0")
			if tc.wantErr == "" {
				if err != nil || string(got) != string(payload) {
					t.Fatalf("got %q, %v", got, err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got err %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestNewSetupOptions(t *testing.T) {
	old := isTerminal
	t.Cleanup(func() { isTerminal = old })

	isTerminal = func(*os.File) bool { return true }
	t.Setenv(nonInteractiveEnv, "")
	if o := newSetupOptions(false, false); !o.prompt || !o.skills {
		t.Errorf("terminal: %+v", o)
	}
	if o := newSetupOptions(true, true); o.prompt || !o.assumeYes || o.skills {
		t.Errorf("--yes --no-skills: %+v", o)
	}
	t.Setenv(nonInteractiveEnv, "1")
	if o := newSetupOptions(false, false); o.prompt {
		t.Errorf("%s=1 still prompts", nonInteractiveEnv)
	}

	t.Setenv(nonInteractiveEnv, "")
	isTerminal = func(*os.File) bool { return false }
	if o := newSetupOptions(false, false); o.prompt || o.assumeYes {
		t.Errorf("no terminal: %+v", o)
	}
}
