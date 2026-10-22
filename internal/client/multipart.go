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
	parts := make([]MultipartPart, 0, len(fields)+1)
	for k, v := range fields {
		parts = append(parts, MultipartPart{FieldName: k, Data: []byte(v)})
	}
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	parts = append(parts, MultipartPart{FieldName: fieldName, FileName: name, ContentType: contentType, Data: data})
	if err := r.SetMultipart(parts); err != nil {
		return err
	}
	summary, _ := json.Marshal(struct {
		Fields map[string]string `json:"fields"`
		File   string            `json:"file"`
		Size   int               `json:"size"`
	}{fields, path, len(data)})
	r.bodySummary = string(summary)
	return nil
}

// MultipartPart is one part of a multipart/form-data body. A part with a
// FileName is a file part; diagnostics show its size, never its Data.
type MultipartPart struct {
	FieldName   string
	FileName    string
	ContentType string // omitted from the part header when empty
	Data        []byte
}

// multipartPartSummary is how a part appears in dry-run and debug output.
type multipartPartSummary struct {
	Field       string `json:"field"`
	File        string `json:"file,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        *int   `json:"size,omitempty"`
	Value       string `json:"value,omitempty"`
}

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// SetMultipart buffers parts as the request body, so the authentication and
// 429 retry paths can replay it.
func (r *Request) SetMultipart(parts []MultipartPart) error {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	summary := make([]multipartPartSummary, 0, len(parts))
	for _, p := range parts {
		if p.FieldName == "" || strings.ContainsAny(p.FieldName+p.FileName+p.ContentType, "\r\n") {
			return fmt.Errorf("invalid multipart part %q: names and content type must be non-empty single lines", p.FieldName)
		}
		header := make(textproto.MIMEHeader)
		s := multipartPartSummary{Field: p.FieldName, ContentType: p.ContentType}
		if p.FileName != "" {
			// Quoted UTF-8 with \ and " escaped, as CreateFormFile writes. Unlike
			// mime.FormatMediaType it never emits filename*=, which RFC 7578 forbids.
			header.Set("Content-Disposition", multipart.FileContentDisposition(p.FieldName, p.FileName))
			size := len(p.Data)
			s.File, s.Size = p.FileName, &size
		} else {
			header.Set("Content-Disposition", `form-data; name="`+quoteEscaper.Replace(p.FieldName)+`"`)
			s.Value = string(p.Data)
		}
		if p.ContentType != "" {
			header.Set("Content-Type", p.ContentType)
		}
		part, err := writer.CreatePart(header)
		if err != nil {
			return fmt.Errorf("creating multipart part %q: %w", p.FieldName, err)
		}
		if _, err := part.Write(p.Data); err != nil {
			return fmt.Errorf("writing multipart part %q: %w", p.FieldName, err)
		}
		summary = append(summary, s)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing multipart body: %w", err)
	}
	summaryJSON, _ := json.Marshal(summary)
	r.bodyRaw = buf.String()
	r.bodyContentType = writer.FormDataContentType()
	r.bodySummary = string(summaryJSON)
	return nil
}
