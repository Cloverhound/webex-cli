package calling

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/audio"
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
	buildURL func(orgID, locationID, announcementID string) string,
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
			f, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening audio file: %w", err)
			}
			defer f.Close()

			url := buildURL(orgID, locationID, announcementID)

			parts := []audio.MultipartPart{
				{
					FieldName:   "name",
					ContentType: "text/plain",
					Data:        strings.NewReader(name),
				},
				{
					FieldName:   "file",
					FileName:    filepath.Base(filePath),
					ContentType: "audio/wav",
					Data:        f,
				},
			}

			body, statusCode, err := audio.UploadMultipart(method, url, parts)
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
		func(orgID, _, _ string) string {
			u := config.CallingBaseURL + "/telephony/config/announcements"
			if orgID != "" {
				u += "?orgId=" + url.QueryEscape(orgID)
			}
			return u
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
		func(orgID, locationID, _ string) string {
			u := config.CallingBaseURL + "/telephony/config/locations/" + url.PathEscape(locationID) + "/announcements"
			if orgID != "" {
				u += "?orgId=" + url.QueryEscape(orgID)
			}
			return u
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
		func(orgID, _, announcementID string) string {
			u := config.CallingBaseURL + "/telephony/config/announcements/" + url.PathEscape(announcementID)
			if orgID != "" {
				u += "?orgId=" + url.QueryEscape(orgID)
			}
			return u
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
		func(orgID, locationID, announcementID string) string {
			u := config.CallingBaseURL + "/telephony/config/locations/" + url.PathEscape(locationID) + "/announcements/" + url.PathEscape(announcementID)
			if orgID != "" {
				u += "?orgId=" + url.QueryEscape(orgID)
			}
			return u
		},
		true, true,
	)
}
