package apicmd

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexops"
	"github.com/spf13/cobra"
)

type formPart struct {
	name, file, contentType, data string
}

func readForm(t *testing.T, body []byte, contentType string) []formPart {
	t.Helper()
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "multipart/form-data" {
		t.Fatalf("content type %q (err %v)", contentType, err)
	}
	var parts []formPart
	r := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	for {
		p, err := r.NextPart()
		if err == io.EOF {
			return parts
		}
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(p)
		parts = append(parts, formPart{p.FormName(), p.FileName(), p.Header.Get("Content-Type"), string(data)})
	}
}

func uploadValues(flags map[string]string) *flagValues {
	values := &flagValues{strings: map[string]*string{}}
	for k, v := range flags {
		values.strings[k] = &v
	}
	return values
}

var testUpload = &webexops.Upload{Parts: []webexops.Part{
	{Name: "info", Kind: webexops.PartJSON, Required: true},
	{Name: "image", Kind: webexops.PartFile, Flag: "file"},
	{Name: "imageName", Kind: webexops.PartFileName, Flag: "file-name"},
}}

func TestUploadFormSendsEachPartInOrder(t *testing.T) {
	file := filepath.Join(t.TempDir(), "logo.png")
	if err := os.WriteFile(file, []byte("\x89PNG image"), 0o600); err != nil {
		t.Fatal(err)
	}
	form, err := uploadForm(testUpload, uploadValues(map[string]string{"file": file, "file-name": ""}), `{"name":"x"}`, "")
	if err != nil {
		t.Fatal(err)
	}
	got := readForm(t, form.Body, form.ContentType)
	want := []formPart{
		{"info", "", "application/json", `{"name":"x"}`},
		{"image", "logo.png", "image/png", "\x89PNG image"},
		{"imageName", "", "", "logo.png"},
	}
	if len(got) != len(want) {
		t.Fatalf("parts %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("part %d is %q, want %q", i, got[i], want[i])
		}
	}
	if strings.Contains(form.Summary, "PNG image") || !strings.Contains(form.Summary, `"file":"logo.png"`) {
		t.Errorf("summary %s should name the file without its contents", form.Summary)
	}
}

func TestUploadFormReadsItsJSONFromTheBodyFile(t *testing.T) {
	body := filepath.Join(t.TempDir(), "channel.json")
	if err := os.WriteFile(body, []byte(`{"from":"file"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	form, err := uploadForm(testUpload, uploadValues(map[string]string{"file": "", "file-name": ""}), `{"from":"flag"}`, body)
	if err != nil {
		t.Fatal(err)
	}
	got := readForm(t, form.Body, form.ContentType)
	if len(got) != 1 || got[0].data != `{"from":"file"}` {
		t.Errorf("parts %q: want only the body file, and no file name without a file", got)
	}
}

func TestUploadFormGivesTheNameItIsGiven(t *testing.T) {
	file := filepath.Join(t.TempDir(), "logo.png")
	if err := os.WriteFile(file, []byte("png"), 0o600); err != nil {
		t.Fatal(err)
	}
	form, err := uploadForm(testUpload, uploadValues(map[string]string{"file": file, "file-name": "Lobby"}), "{}", "")
	if err != nil {
		t.Fatal(err)
	}
	if got := readForm(t, form.Body, form.ContentType); len(got) != 3 || got[2].data != "Lobby" {
		t.Errorf("parts %q, want the name Lobby", got)
	}
}

func TestUploadFormRefusesAMissingRequiredBody(t *testing.T) {
	_, err := uploadForm(testUpload, uploadValues(map[string]string{"file": "", "file-name": ""}), "", "")
	if err == nil || !strings.Contains(err.Error(), "--body or --body-file is required") {
		t.Errorf("err = %v", err)
	}
}

// `--file "$UNSET"` reaches the command as an empty path, which must not send a
// form without its file.
func TestUploadFormRefusesAnEmptyPathForARequiredFile(t *testing.T) {
	upload := &webexops.Upload{Parts: []webexops.Part{{Name: "file", Kind: webexops.PartFile, Flag: "file", Required: true}}}
	_, err := uploadForm(upload, uploadValues(map[string]string{"file": ""}), "", "")
	if err == nil || !strings.Contains(err.Error(), "--file needs a path") {
		t.Errorf("err = %v", err)
	}
}

func TestUploadFormReportsAFileItCannotRead(t *testing.T) {
	_, err := uploadForm(testUpload, uploadValues(map[string]string{"file": filepath.Join(t.TempDir(), "missing.png"), "file-name": ""}), "{}", "")
	if err == nil || !strings.Contains(err.Error(), "reading --file") {
		t.Errorf("err = %v", err)
	}
}

// Every upload in the table gets a flag for each part, and no flags for the
// collection's body.
func TestUploadCommandsTakeEachPartAsAFlag(t *testing.T) {
	uploads := 0
	for _, area := range webexops.Areas() {
		for _, op := range webexops.Operations(area) {
			if op.Upload == nil {
				continue
			}
			uploads++
			if op.Body != nil {
				t.Errorf("%s has both a body and a form", op.ID)
			}
			c := NewCommand(op)
			for _, part := range op.Upload.Parts {
				flag := part.Flag
				if part.Kind == webexops.PartJSON {
					flag = "body"
					if c.Flags().Lookup("body-file") == nil {
						t.Errorf("%s: no --body-file for the part %s", op.ID, part.Name)
					}
				}
				f := c.Flags().Lookup(flag)
				switch {
				case f == nil:
					t.Errorf("%s: no --%s for the part %s", op.ID, flag, part.Name)
				case part.Kind != webexops.PartJSON && part.Required != (len(f.Annotations[cobra.BashCompOneRequiredFlag]) > 0):
					t.Errorf("%s: --%s required = %v, want %v", op.ID, flag, !part.Required, part.Required)
				}
			}
		}
	}
	if uploads != 6 {
		t.Errorf("%d uploads in the table, want 6", uploads)
	}
}
