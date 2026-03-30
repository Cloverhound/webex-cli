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
	registerAudioFileDownload()
	registerAudioFileUpload()
}

func registerAudioFileDownload() {
	var orgid string
	var id string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download an audio file to disk",
		Long: `Download a Contact Center audio file binary to a local file.

Fetches the audio file metadata (with download URL), then downloads the
binary content and writes it to the specified output path.

Examples:
  webex cc audio-files download --id <audio-file-id> --output prompt.wav
  webex cc audio-files download --id <audio-file-id> --output /tmp/greeting.wav --debug`,
		RunE: func(c *cobra.Command, args []string) error {
			// Step 1: GET metadata with includeUrl=true
			req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/audio-file/{id}")
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
				return fmt.Errorf("parsing audio file metadata: %w", err)
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
	cmd.Flags().StringVar(&id, "id", "", "Audio file resource ID")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVar(&outputPath, "output", "", "File path to write audio to")
	cmd.MarkFlagRequired("output")

	audioFilesCmd.AddCommand(cmd)
}

func registerAudioFileUpload() {
	var orgid string
	var filePath string
	var name string
	var description string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a WAV file as a new audio file",
		Long: `Upload an audio file (WAV) to Contact Center using multipart/form-data.

The upload sends two parts: audioFileInfo (JSON metadata) and audioFile (binary).
If --name is omitted, the filename (with .wav extension) is used.

Examples:
  webex cc audio-files upload --file prompt.wav
  webex cc audio-files upload --file prompt.wav --name "Main Greeting"
  webex cc audio-files upload --file prompt.wav --name "Main Greeting" --description "IVR main menu"
  webex cc audio-files upload --file prompt.wav --dry-run
  webex cc audio-files upload --file prompt.wav --debug`,
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

			u := config.CcBaseURL + "/organization/" + url.PathEscape(orgid) + "/audio-file"

			info := map[string]any{
				"name":          name,
				"contentType":   "AUDIO_WAV",
				"systemDefault": false,
			}
			if description != "" {
				info["description"] = description
			}
			infoJSON, _ := json.Marshal(info)

			parts := []audio.MultipartPart{
				{
					FieldName:   "audioFileInfo",
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
	cmd.Flags().StringVar(&name, "name", "", "Audio file name (defaults to filename without extension)")
	cmd.Flags().StringVar(&description, "description", "", "Audio file description")

	audioFilesCmd.AddCommand(cmd)
}
