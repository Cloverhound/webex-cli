package audio

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

// MultipartPart describes one part of a multipart/form-data upload.
type MultipartPart struct {
	FieldName   string    // form field name
	FileName    string    // original filename (empty for non-file parts)
	ContentType string    // e.g. "application/json", "audio/wav"
	Data        io.Reader // part content
}

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

// UploadMultipart builds a multipart/form-data request and sends it.
// Returns response body, HTTP status code, and error.
func UploadMultipart(method, url string, parts []MultipartPart) ([]byte, int, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, p := range parts {
		header := make(textproto.MIMEHeader)
		if p.FileName != "" {
			header.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`, p.FieldName, p.FileName))
		} else {
			header.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"`, p.FieldName))
		}
		header.Set("Content-Type", p.ContentType)

		part, err := writer.CreatePart(header)
		if err != nil {
			return nil, 0, fmt.Errorf("creating multipart part %q: %w", p.FieldName, err)
		}
		if _, err := io.Copy(part, p.Data); err != nil {
			return nil, 0, fmt.Errorf("writing multipart part %q: %w", p.FieldName, err)
		}
	}
	writer.Close()

	contentType := writer.FormDataContentType()

	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", method, url)
		fmt.Fprintf(os.Stderr, "DEBUG:   Content-Type: %s\n", contentType)
		for _, p := range parts {
			if p.FileName != "" {
				fmt.Fprintf(os.Stderr, "DEBUG:   Part %q: file=%s (%d bytes total body)\n",
					p.FieldName, p.FileName, buf.Len())
			} else {
				fmt.Fprintf(os.Stderr, "DEBUG:   Part %q: %s\n", p.FieldName, p.ContentType)
			}
		}
	}

	if config.DryRun() {
		fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", method, url)
		fmt.Fprintf(os.Stderr, "[DRY RUN]   Content-Type: %s\n", contentType)
		for _, p := range parts {
			if p.FileName != "" {
				fmt.Fprintf(os.Stderr, "[DRY RUN]   Part %q: file=%s\n", p.FieldName, p.FileName)
			} else {
				fmt.Fprintf(os.Stderr, "[DRY RUN]   Part %q\n", p.FieldName)
			}
		}
		return nil, 0, client.ErrDryRun
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, 0, fmt.Errorf("building upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token())
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("uploading: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading upload response: %w", err)
	}

	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", resp.StatusCode, len(body))
	}

	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, resp.StatusCode, nil
}
