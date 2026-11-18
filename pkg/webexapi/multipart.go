package webexapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// MultipartPart is one part of a multipart/form-data body. A part with a
// FileName is a file part; Summary shows its size, never its Data.
type MultipartPart struct {
	FieldName   string
	FileName    string
	ContentType string // left out of the part header when empty
	Data        []byte
}

// Multipart is an encoded multipart/form-data body. It is held in memory, so a
// retry after a 401 or 429 can send it again.
type Multipart struct {
	Body        []byte
	ContentType string

	// Summary describes the parts as JSON for diagnostics: each field's value,
	// and each file's name, type and size.
	Summary string
}

type partSummary struct {
	Field       string `json:"field"`
	File        string `json:"file,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        *int   `json:"size,omitempty"`
	Value       string `json:"value,omitempty"`
}

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

// EncodeMultipart encodes parts as a multipart/form-data body.
func EncodeMultipart(parts []MultipartPart) (Multipart, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	summary := make([]partSummary, 0, len(parts))
	for _, p := range parts {
		if p.FieldName == "" || strings.ContainsAny(p.FieldName+p.FileName+p.ContentType, "\r\n") {
			return Multipart{}, fmt.Errorf("invalid multipart part %q: names and content type must be non-empty single lines", p.FieldName)
		}
		header := make(textproto.MIMEHeader)
		s := partSummary{Field: p.FieldName, ContentType: p.ContentType}
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
			return Multipart{}, fmt.Errorf("creating multipart part %q: %w", p.FieldName, err)
		}
		if _, err := part.Write(p.Data); err != nil {
			return Multipart{}, fmt.Errorf("writing multipart part %q: %w", p.FieldName, err)
		}
		summary = append(summary, s)
	}
	if err := writer.Close(); err != nil {
		return Multipart{}, fmt.Errorf("closing multipart body: %w", err)
	}
	summaryJSON, _ := json.Marshal(summary)
	return Multipart{Body: buf.Bytes(), ContentType: writer.FormDataContentType(), Summary: string(summaryJSON)}, nil
}

// MultipartFile encodes string form fields, in name order, followed by one file
// part read from path. The file must be a regular file of at most limit bytes; its
// content type comes from its extension, or from its content.
func MultipartFile(fields map[string]string, fieldName, path string, limit int64) (Multipart, error) {
	if fieldName == "" || strings.ContainsAny(fieldName, "\r\n") {
		return Multipart{}, fmt.Errorf("invalid attachment field name")
	}
	for name := range fields {
		if name == "" || strings.ContainsAny(name, "\r\n") {
			return Multipart{}, fmt.Errorf("invalid multipart field name %q", name)
		}
	}
	tooLarge := fmt.Errorf("attachment %q exceeds the %s limit", path, sizeLabel(limit))
	check := func(info os.FileInfo, err error) error {
		switch {
		case err != nil:
			return fmt.Errorf("checking attachment %q: %w", path, err)
		case !info.Mode().IsRegular():
			return fmt.Errorf("attachment %q must be a regular file", path)
		case info.Size() > limit:
			return tooLarge
		}
		return nil
	}
	// Check the path before opening so FIFOs and devices are refused without
	// blocking, then check the opened descriptor in case the path changed.
	if err := check(os.Stat(path)); err != nil {
		return Multipart{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return Multipart{}, fmt.Errorf("opening attachment %q: %w", path, err)
	}
	defer f.Close()
	if err := check(f.Stat()); err != nil {
		return Multipart{}, err
	}
	name := filepath.Base(path)
	if strings.ContainsAny(name, "\r\n") {
		return Multipart{}, fmt.Errorf("attachment filename must not contain newlines")
	}
	// Reading one byte past the limit shows a file that grew past it after Stat.
	var r io.Reader = f
	if limit < math.MaxInt64 {
		r = io.LimitReader(f, limit+1)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return Multipart{}, fmt.Errorf("reading attachment %q: %w", path, err)
	}
	if int64(len(data)) > limit {
		return Multipart{}, tooLarge
	}
	parts := make([]MultipartPart, 0, len(fields)+1)
	for _, k := range slices.Sorted(maps.Keys(fields)) {
		parts = append(parts, MultipartPart{FieldName: k, Data: []byte(fields[k])})
	}
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	parts = append(parts, MultipartPart{FieldName: fieldName, FileName: name, ContentType: contentType, Data: data})
	return EncodeMultipart(parts)
}

func sizeLabel(limit int64) string {
	if limit%(1<<20) == 0 {
		return fmt.Sprintf("%d MiB", limit>>20)
	}
	return fmt.Sprintf("%d-byte", limit)
}
