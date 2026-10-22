package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// MaxMessageFileSize is the local message attachment limit (100 MiB).
const MaxMessageFileSize int64 = 100 * 1024 * 1024

// SetMultipartFile buffers one bounded regular file and string form fields.
// The buffered body is replayable by the normal authentication/429 retry path.
// Diagnostics identify the attachment without logging its contents.
func (r *Request) SetMultipartFile(fields map[string]string, fieldName, path string) error {
	if fieldName == "" || strings.ContainsAny(fieldName, "\r\n") {
		return fmt.Errorf("invalid attachment field name")
	}
	for name := range fields {
		if name == "" || strings.ContainsAny(name, "\r\n") {
			return fmt.Errorf("invalid multipart field name %q", name)
		}
	}
	tooLarge := fmt.Errorf("attachment %q exceeds the 100 MiB limit", path)
	check := func(info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return fmt.Errorf("checking attachment %q: %w", path, err)
		case !info.Mode().IsRegular():
			return fmt.Errorf("attachment %q must be a regular file", path)
		case info.Size() > MaxMessageFileSize:
			return tooLarge
		}
		return nil
	}
	// Check the path before opening so FIFOs/devices are refused without
	// blocking, then check the opened descriptor in case the path changed.
	if err := check(os.Stat(path)); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening attachment %q: %w", path, err)
	}
	defer f.Close()
	if err := check(f.Stat()); err != nil {
		return err
	}
	name := filepath.Base(path)
	if strings.ContainsAny(name, "\r\n") {
		return fmt.Errorf("attachment filename must not contain newlines")
	}
	data, err := io.ReadAll(io.LimitReader(f, MaxMessageFileSize+1))
	if err != nil {
		return fmt.Errorf("reading attachment %q: %w", path, err)
	}
	if int64(len(data)) > MaxMessageFileSize {
		return tooLarge
	}
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return fmt.Errorf("writing message field %q: %w", k, err)
		}
	}
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	header := make(textproto.MIMEHeader)
	// Quoted UTF-8 with \ and " escaped, as CreateFormFile writes. Unlike
	// mime.FormatMediaType it never emits filename*=, which RFC 7578 forbids.
	header.Set("Content-Disposition", multipart.FileContentDisposition(fieldName, name))
	header.Set("Content-Type", contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return fmt.Errorf("creating attachment part: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return fmt.Errorf("writing attachment: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing attachment body: %w", err)
	}
	summary, _ := json.Marshal(struct {
		Fields map[string]string `json:"fields"`
		File   string            `json:"file"`
		Size   int               `json:"size"`
	}{fields, path, len(data)})
	r.bodyRaw = buf.String()
	r.bodyContentType = writer.FormDataContentType()
	r.bodySummary = string(summary)
	return nil
}
