package cc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

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

			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Downloading from %s\n", dlURL)
			}

			// Step 3: Download the binary (pre-signed S3 URL, no auth header needed)
			dlReq, err := http.NewRequest("GET", dlURL, nil)
			if err != nil {
				return fmt.Errorf("building download request: %w", err)
			}

			dlResp, err := http.DefaultClient.Do(dlReq)
			if err != nil {
				return fmt.Errorf("downloading audio file: %w", err)
			}
			defer dlResp.Body.Close()

			if dlResp.StatusCode >= 400 {
				body, _ := io.ReadAll(dlResp.Body)
				return fmt.Errorf("download failed with status %d: %s", dlResp.StatusCode, string(body))
			}

			// Step 4: Write to output file
			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()

			n, err := io.Copy(f, dlResp.Body)
			if err != nil {
				return fmt.Errorf("writing audio file: %w", err)
			}

			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Wrote %d bytes to %s\n", n, outputPath)
			}

			// Step 5: Print metadata to stdout
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
			// Open and validate the file
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

			// Build the URL
			url := config.CcBaseURL + "/organization/" + orgid + "/audio-file"

			// Build multipart body
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			// Part 1: audioFileInfo JSON
			infoHeader := make(textproto.MIMEHeader)
			infoHeader.Set("Content-Disposition", `form-data; name="audioFileInfo"`)
			infoHeader.Set("Content-Type", "application/json")
			infoPart, err := writer.CreatePart(infoHeader)
			if err != nil {
				return fmt.Errorf("creating audioFileInfo part: %w", err)
			}

			info := map[string]any{
				"name":          name,
				"contentType":   "AUDIO_WAV",
				"systemDefault": false,
			}
			if description != "" {
				info["description"] = description
			}
			infoJSON, _ := json.Marshal(info)
			infoPart.Write(infoJSON)

			// Part 2: audioFile binary
			fileHeader := make(textproto.MIMEHeader)
			fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="audioFile"; filename="%s"`, filepath.Base(filePath)))
			fileHeader.Set("Content-Type", "audio/wav")
			filePart, err := writer.CreatePart(fileHeader)
			if err != nil {
				return fmt.Errorf("creating audioFile part: %w", err)
			}
			if _, err := io.Copy(filePart, f); err != nil {
				return fmt.Errorf("writing audio file part: %w", err)
			}

			writer.Close()

			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: POST %s\n", url)
				fmt.Fprintf(os.Stderr, "DEBUG:   Content-Type: %s\n", writer.FormDataContentType())
				fmt.Fprintf(os.Stderr, "DEBUG:   audioFileInfo: %s\n", string(infoJSON))
				fmt.Fprintf(os.Stderr, "DEBUG:   audioFile: %s (%d bytes)\n", filepath.Base(filePath), buf.Len())
			}

			// Dry-run: print what would be sent and return
			if config.DryRun() {
				fmt.Fprintf(os.Stderr, "[DRY RUN] POST %s\n", url)
				fmt.Fprintf(os.Stderr, "[DRY RUN]   Content-Type: %s\n", writer.FormDataContentType())
				fmt.Fprintf(os.Stderr, "[DRY RUN]   audioFileInfo: %s\n", string(infoJSON))
				fmt.Fprintf(os.Stderr, "[DRY RUN]   audioFile: %s\n", filepath.Base(filePath))
				return client.ErrDryRun
			}

			// Execute the request
			req, err := http.NewRequest("POST", url, &buf)
			if err != nil {
				return fmt.Errorf("building upload request: %w", err)
			}
			req.Header.Set("Authorization", "Bearer "+config.Token())
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("uploading audio file: %w", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("reading upload response: %w", err)
			}

			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", resp.StatusCode, len(body))
			}

			if resp.StatusCode >= 400 {
				return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
			}

			return output.Print(body, resp.StatusCode)
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
