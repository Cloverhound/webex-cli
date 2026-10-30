package audio

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// DownloadBinary fetches a URL (typically a pre-signed S3 link) and writes
// the response body to outputPath. Returns bytes written.
func DownloadBinary(url, outputPath string) (int64, error) {
	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: Downloading from %s\n", url)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("building download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		return 0, fmt.Errorf("writing file: %w", err)
	}

	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: Wrote %d bytes to %s\n", n, outputPath)
	}

	return n, nil
}
