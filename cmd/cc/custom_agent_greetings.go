package cc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/audio"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	registerAgentGreetingDownload()
	registerAgentGreetingUpload()
}

func registerAgentGreetingDownload() {
	var orgid string
	var id string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download an agent personal greeting file to disk",
		Long: `Download a Contact Center agent personal greeting file to a local file.

Fetches the greeting file metadata (with download URL), then downloads the
binary content and writes it to the specified output path.

Examples:
  webex cc agent-personal-greeting-files download --orgid <org-id> --id <greeting-id> --output greeting.wav
  webex cc agent-personal-greeting-files download --orgid <org-id> --id <greeting-id> --output /tmp/greeting.wav --debug`,
		RunE: func(c *cobra.Command, args []string) error {
			// Step 1: GET metadata with includeUrl=true
			req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/agent-personal-greeting/{id}")
			req.PathParam("orgid", orgid)
			req.PathParam("id", id)
			req.QueryParam("includeUrl", "true")

			resp, statusCode, err := req.Do()
			if err != nil {
				return err
			}

			// Step 2: Extract download URL from response
			var meta map[string]any
			if err := json.Unmarshal(resp, &meta); err != nil {
				return fmt.Errorf("parsing greeting metadata: %w", err)
			}

			dlURL, ok := meta["url"].(string)
			if !ok || dlURL == "" {
				return fmt.Errorf("no download URL in response (status %d)", statusCode)
			}

			// Step 3: Download the binary
			if _, err := audio.DownloadBinary(dlURL, outputPath); err != nil {
				return err
			}

			// Step 4: Print metadata to stdout
			return output.Print(resp, statusCode)
		},
	}

	cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID")
	cmd.MarkFlagRequired("orgid")
	cmd.Flags().StringVar(&id, "id", "", "Agent personal greeting file ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&outputPath, "output", "", "File path to write greeting to")
	cmd.MarkFlagRequired("output")

	agentPersonalGreetingFilesCmd.AddCommand(cmd)
}

func registerAgentGreetingUpload() {
	var orgid string
	var filePath string
	var name string
	var agentID string
	var greetingPurposeID string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a WAV file as a new agent personal greeting",
		Long: `Upload an agent personal greeting file (WAV) to Contact Center using multipart/form-data.

The upload sends two parts: agentPersonalGreetingInfo (JSON metadata) and audioFile (binary).
If --name is omitted, the filename (with .wav extension) is used.

Examples:
  webex cc agent-personal-greeting-files upload --agent-id <agent-id> --file greeting.wav
  webex cc agent-personal-greeting-files upload --agent-id <agent-id> --file greeting.wav --name "Personal Greeting"
  webex cc agent-personal-greeting-files upload --agent-id <agent-id> --file greeting.wav --greeting-purpose-id <id>
  webex cc agent-personal-greeting-files upload --agent-id <agent-id> --file greeting.wav --dry-run`,
		RunE: func(c *cobra.Command, args []string) error {
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening audio file: %w", err)
			}
			defer f.Close()

			// Default name to filename; ensure it ends with .wav (API requires it)
			if name == "" {
				name = filepath.Base(filePath)
			}
			if !strings.HasSuffix(strings.ToLower(name), ".wav") {
				name += ".wav"
			}

			u := config.CcBaseURL + "/organization/" + url.PathEscape(orgid) + "/v2/agent-personal-greeting"

			info := map[string]any{
				"agentId":     agentID,
				"name":        name,
				"contentType": "AUDIO_WAV",
			}
			if greetingPurposeID != "" {
				info["greetingPurposeId"] = greetingPurposeID
			}
			infoJSON, _ := json.Marshal(info)

			parts := []audio.MultipartPart{
				{
					FieldName:   "agentPersonalGreetingInfo",
					ContentType: "application/json",
					Data:        strings.NewReader(string(infoJSON)),
				},
				{
					FieldName:   "audioFile",
					FileName:    filepath.Base(filePath),
					ContentType: "audio/wav",
					Data:        f,
				},
			}

			body, statusCode, err := audio.UploadMultipart("POST", u, parts)
			if err != nil {
				return err
			}

			return output.Print(body, statusCode)
		},
	}

	cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID")
	cmd.MarkFlagRequired("orgid")
	cmd.Flags().StringVar(&filePath, "file", "", "Path to WAV file to upload")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&agentID, "agent-id", "", "Agent ID (WxCC agent resource ID)")
	cmd.MarkFlagRequired("agent-id")
	cmd.Flags().StringVar(&name, "name", "", "Greeting name (defaults to filename)")
	cmd.Flags().StringVar(&greetingPurposeID, "greeting-purpose-id", "", "Greeting purpose ID")

	agentPersonalGreetingFilesCmd.AddCommand(cmd)
}
