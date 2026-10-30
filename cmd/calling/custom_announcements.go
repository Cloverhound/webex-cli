package calling

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	registerAnnouncementUpload()
	registerAnnouncementUploadLocation()
	registerAnnouncementUpdate()
	registerAnnouncementUpdateLocation()
}

// registerAnnouncementCmd is a helper that registers one announcement upload/update command.
func registerAnnouncementCmd(use, short, long string, method string,
	buildPath func(locationID, announcementID string) string,
	needsLocation, needsAnnouncement bool) {

	var filePath string
	var name string
	var orgID string
	var locationID string
	var announcementID string

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		RunE: func(c *cobra.Command, args []string) error {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("reading audio file: %w", err)
			}

			req := client.NewRequest(config.CallingBaseURL, method, buildPath(locationID, announcementID))
			req.QueryParam("orgId", orgID)
			if err := req.SetMultipart([]client.MultipartPart{
				{
					FieldName:   "name",
					ContentType: "text/plain",
					Data:        []byte(name),
				},
				{
					FieldName:   "file",
					FileName:    filepath.Base(filePath),
					ContentType: "audio/wav",
					Data:        data,
				},
			}); err != nil {
				return err
			}

			body, statusCode, err := req.Do()
			if err != nil {
				return err
			}

			return output.Print(body, statusCode)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to WAV file to upload")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVar(&name, "name", "", "Announcement name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	if needsLocation {
		cmd.Flags().StringVar(&locationID, "location-id", "", "Location ID")
		cmd.MarkFlagRequired("location-id")
	}
	if needsAnnouncement {
		cmd.Flags().StringVar(&announcementID, "announcement-id", "", "Announcement ID")
		cmd.MarkFlagRequired("announcement-id")
	}

	announcementRepositoryCmd.AddCommand(cmd)
}

func registerAnnouncementUpload() {
	registerAnnouncementCmd(
		"upload-binary-greeting",
		"Upload Binary Announcement Greeting",
		`Upload an announcement greeting (WAV) at the organization level.

Uses multipart/form-data with a "name" text field and a "file" binary part.

Examples:
  webex calling announcement-repository upload-binary-greeting --file greeting.wav --name "Main Greeting"
  webex calling announcement-repository upload-binary-greeting --file greeting.wav --name "Main Greeting" --dry-run`,
		"POST",
		func(_, _ string) string {
			return "/telephony/config/announcements"
		},
		false, false,
	)
}

func registerAnnouncementUploadLocation() {
	registerAnnouncementCmd(
		"upload-binary-greeting-2",
		"Upload Binary Announcement Greeting at Location",
		`Upload an announcement greeting (WAV) at a specific location.

Uses multipart/form-data with a "name" text field and a "file" binary part.

Examples:
  webex calling announcement-repository upload-binary-greeting-2 --file greeting.wav --name "Lobby Greeting" --location-id <loc-id>
  webex calling announcement-repository upload-binary-greeting-2 --file greeting.wav --name "Lobby Greeting" --location-id <loc-id> --dry-run`,
		"POST",
		func(locationID, _ string) string {
			return "/telephony/config/locations/" + url.PathEscape(locationID) + "/announcements"
		},
		true, false,
	)
}

func registerAnnouncementUpdate() {
	registerAnnouncementCmd(
		"update-binary-greeting",
		"Update Binary Announcement Greeting",
		`Update an existing announcement greeting (WAV) at the organization level.

Uses multipart/form-data with a "name" text field and a "file" binary part.

Examples:
  webex calling announcement-repository update-binary-greeting --file greeting.wav --name "Updated Greeting" --announcement-id <ann-id>
  webex calling announcement-repository update-binary-greeting --file greeting.wav --name "Updated Greeting" --announcement-id <ann-id> --dry-run`,
		"PUT",
		func(_, announcementID string) string {
			return "/telephony/config/announcements/" + url.PathEscape(announcementID)
		},
		false, true,
	)
}

func registerAnnouncementUpdateLocation() {
	registerAnnouncementCmd(
		"update-binary-greeting-2",
		"Update Binary Announcement Greeting at Location",
		`Update an existing announcement greeting (WAV) at a specific location.

Uses multipart/form-data with a "name" text field and a "file" binary part.

Examples:
  webex calling announcement-repository update-binary-greeting-2 --file greeting.wav --name "Updated Greeting" --location-id <loc-id> --announcement-id <ann-id>`,
		"PUT",
		func(locationID, announcementID string) string {
			return "/telephony/config/locations/" + url.PathEscape(locationID) + "/announcements/" + url.PathEscape(announcementID)
		},
		true, true,
	)
}
