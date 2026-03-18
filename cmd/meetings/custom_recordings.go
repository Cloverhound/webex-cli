package meetings

import (
	"encoding/json"
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/audio"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	registerRecordingsDownload()
}

func registerRecordingsDownload() {
	var recordingID string
	var outputPath string
	var dlType string
	var hostEmail string

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a meeting recording to disk",
		Long: `Download a meeting recording, audio, or transcript to a local file.

Fetches the recording metadata to obtain the temporary download link,
then downloads the binary content and writes it to the specified output path.

The --type flag selects which download link to use:
  audio      - MP3 audio track (default)
  recording  - Full video recording (MP4)
  transcript - VTT transcript file

Examples:
  webex meetings recordings download --recording-id <id> --output meeting.mp3
  webex meetings recordings download --recording-id <id> --output meeting.mp4 --type recording
  webex meetings recordings download --recording-id <id> --output meeting.vtt --type transcript
  webex meetings recordings download --recording-id <id> --output meeting.mp3 --host-email admin@example.com`,
		RunE: func(c *cobra.Command, args []string) error {
			// Step 1: GET recording metadata
			req := client.NewRequest(config.CallingBaseURL, "GET", "/recordings/{recordingId}")
			req.PathParam("recordingId", recordingID)
			req.QueryParam("hostEmail", hostEmail)

			resp, statusCode, err := req.Do()
			if err != nil {
				return err
			}

			// Step 2: Extract download URL from temporaryDirectDownloadLinks
			var meta map[string]any
			if err := json.Unmarshal(resp, &meta); err != nil {
				return fmt.Errorf("parsing recording metadata: %w", err)
			}

			links, ok := meta["temporaryDirectDownloadLinks"].(map[string]any)
			if !ok {
				return fmt.Errorf("no temporaryDirectDownloadLinks in response (status %d)", statusCode)
			}

			linkKey := "audioDownloadLink"
			switch dlType {
			case "recording":
				linkKey = "recordingDownloadLink"
			case "transcript":
				linkKey = "transcriptDownloadLink"
			}

			dlURL, ok := links[linkKey].(string)
			if !ok || dlURL == "" {
				return fmt.Errorf("no %s in response", linkKey)
			}

			// Step 3: Download the binary
			if _, err := audio.DownloadBinary(dlURL, outputPath); err != nil {
				return err
			}

			// Step 4: Print metadata to stdout
			return output.Print(resp, statusCode)
		},
	}

	cmd.Flags().StringVar(&recordingID, "recording-id", "", "Recording ID")
	cmd.MarkFlagRequired("recording-id")
	cmd.Flags().StringVar(&outputPath, "output", "", "File path to write recording to")
	cmd.MarkFlagRequired("output")
	cmd.Flags().StringVar(&dlType, "type", "audio", "Download type: audio, recording, or transcript")
	cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email of the recording host (admin access to other users' recordings)")

	recordingsCmd.AddCommand(cmd)
}
