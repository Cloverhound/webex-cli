package messaging

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() { messagesCmd.AddCommand(newMessageCreateCommand()) }

func newMessageCreateCommand() *cobra.Command {
	var bodyRaw, bodyFile string
	c := &cobra.Command{
		Use: "create", Short: "Create a Message",
		Long: `Send a message, optionally with one remote URL or local file in the JSON files array.
Local files are uploaded as multipart/form-data; HTTP(S) URLs retain the JSON request.
Relative attachment paths resolve against the current working directory, not the body file.
Local uploads accept string-valued message fields and one regular file up to 100 MiB.

Examples:
  webex messaging messages create --body '{"toPersonEmail":"user@example.com","files":["./guide.md"]}'
  webex messaging messages create --body-file message.json --dry-run`,
		RunE: func(c *cobra.Command, args []string) error {
			raw := bodyRaw
			if bodyFile != "" {
				b, err := os.ReadFile(bodyFile)
				if err != nil {
					return fmt.Errorf("reading body file: %w", err)
				}
				raw = string(b)
			}
			req, err := messageCreateRequest(raw)
			if err != nil {
				return err
			}
			resp, status, err := req.Do()
			if err != nil {
				return err
			}
			return output.Print(resp, status)
		},
	}
	c.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
	c.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
	return c
}

var attachmentScheme = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*:`)

func messageCreateRequest(raw string) (*client.Request, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &fields); err != nil || fields == nil {
		return nil, fmt.Errorf("message body must be a JSON object")
	}
	req := client.NewRequest(config.CallingBaseURL, "POST", "/messages")
	req.SetBodyRaw(raw)
	value, ok := fields["files"]
	if !ok {
		return req, nil
	}
	var files []string
	if string(value) == "null" {
		return nil, fmt.Errorf("files must be an array of file URLs or local paths")
	}
	if err := json.Unmarshal(value, &files); err != nil {
		return nil, fmt.Errorf("files must be an array of file URLs or local paths: %w", err)
	}
	if len(files) == 0 {
		return req, nil
	}
	if len(files) != 1 {
		return nil, fmt.Errorf("files supports one attachment per message; do not mix local paths and URLs")
	}
	path := files[0]
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("files[0] must not be empty")
	}
	// Drive-letter paths are filesystem paths, not URI schemes. On other platforms
	// they fail with the ordinary local-file error rather than being fetched.
	drivePath := len(path) >= 3 && ((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z')) && path[1] == ':' && (path[2] == '\\' || path[2] == '/')
	if !drivePath && attachmentScheme.MatchString(path) {
		u, err := url.Parse(path)
		if err != nil || u.Hostname() == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return nil, fmt.Errorf("files[0] must be an HTTP(S) URL or local path (file:// is not supported)")
		}
		return req, nil
	}
	values := make(map[string]string, len(fields)-1)
	for name, value := range fields {
		if name == "files" {
			continue
		}
		var text string
		// json.Unmarshal accepts null into a string, so require an actual JSON string.
		if len(value) == 0 || value[0] != '"' || json.Unmarshal(value, &text) != nil {
			return nil, fmt.Errorf("field %q must be a string when uploading a local file; structured fields such as attachments are not supported", name)
		}
		values[name] = text
	}
	if err := req.SetMultipartFile(values, "files", path); err != nil {
		return nil, err
	}
	return req, nil
}
