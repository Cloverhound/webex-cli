package client

import "github.com/Cloverhound/webex-cli/pkg/webexapi"

// MaxMessageFileSize is the local message attachment limit (100 MiB).
const MaxMessageFileSize int64 = 100 * 1024 * 1024

// MultipartPart is one part of a multipart/form-data body; see
// webexapi.MultipartPart.
type MultipartPart = webexapi.MultipartPart

// SetMultipartFile buffers one bounded regular file and string form fields.
// The buffered body is replayable by the normal authentication/429 retry path.
// Diagnostics identify the attachment without logging its contents.
func (r *Request) SetMultipartFile(fields map[string]string, fieldName, path string) error {
	body, err := webexapi.MultipartFile(fields, fieldName, path, MaxMessageFileSize)
	if err != nil {
		return err
	}
	r.setMultipart(body)
	return nil
}

// SetMultipart buffers parts as the request body, so the authentication and
// 429 retry paths can replay it.
func (r *Request) SetMultipart(parts []MultipartPart) error {
	body, err := webexapi.EncodeMultipart(parts)
	if err != nil {
		return err
	}
	r.setMultipart(body)
	return nil
}

func (r *Request) setMultipart(body webexapi.Multipart) {
	r.bodyRaw = ""
	r.bodyBytes = body.Body
	r.bodyContentType = body.ContentType
	r.bodySummary = body.Summary
}
