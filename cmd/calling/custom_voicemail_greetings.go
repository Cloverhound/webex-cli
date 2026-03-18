package calling

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/internal/audio"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

func init() {
	// user-call (person) greetings
	registerGreetingUpload(userCallCmd,
		"update-intercept-greeting-person",
		"Upload Call Intercept Greeting for Person",
		`Upload a call intercept announcement greeting for a person.

Examples:
  webex calling user-call update-intercept-greeting-person --person-id <id> --file greeting.wav
  webex calling user-call update-intercept-greeting-person --person-id <id> --file greeting.wav --dry-run`,
		"/people/{entityId}/features/intercept/actions/announcementUpload/invoke",
		"person-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(userCallCmd,
		"update-busy-voicemail-greeting-person",
		"Upload Busy Voicemail Greeting for Person",
		`Upload a busy voicemail greeting for a person.

Examples:
  webex calling user-call update-busy-voicemail-greeting-person --person-id <id> --file greeting.wav`,
		"/people/{entityId}/features/voicemail/actions/uploadBusyGreeting/invoke",
		"person-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(userCallCmd,
		"update-no-answer-voicemail-greeting-person",
		"Upload No-Answer Voicemail Greeting for Person",
		`Upload a no-answer voicemail greeting for a person.

Examples:
  webex calling user-call update-no-answer-voicemail-greeting-person --person-id <id> --file greeting.wav`,
		"/people/{entityId}/features/voicemail/actions/uploadNoAnswerGreeting/invoke",
		"person-id", func(cmd *cobra.Command) {})

	// call-settings-for-me (self) greetings
	registerGreetingUpload(callSettingsForMeCmd,
		"upload-voicemail-busy-greeting",
		"Upload My Busy Voicemail Greeting",
		`Upload your own busy voicemail greeting.

Examples:
  webex calling call-settings-for-me upload-voicemail-busy-greeting --file greeting.wav
  webex calling call-settings-for-me upload-voicemail-busy-greeting --file greeting.wav --dry-run`,
		"/telephony/config/people/me/settings/voicemail/actions/busyGreetingUpload/invoke",
		"", func(cmd *cobra.Command) {})

	registerGreetingUpload(callSettingsForMeCmd,
		"upload-voicemail-no-answer-greeting",
		"Upload My No-Answer Voicemail Greeting",
		`Upload your own no-answer voicemail greeting.

Examples:
  webex calling call-settings-for-me upload-voicemail-no-answer-greeting --file greeting.wav`,
		"/telephony/config/people/me/settings/voicemail/actions/noAnswerGreetingUpload/invoke",
		"", func(cmd *cobra.Command) {})

	// virtual-line-call greetings
	registerGreetingUpload(virtualLineCallCmd,
		"update-intercept-greeting",
		"Upload Call Intercept Greeting for Virtual Line",
		`Upload a call intercept announcement greeting for a virtual line.

Examples:
  webex calling virtual-line-call update-intercept-greeting --virtual-line-id <id> --file greeting.wav`,
		"/telephony/config/virtualLines/{entityId}/actions/announcementUpload/invoke",
		"virtual-line-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(virtualLineCallCmd,
		"update-busy-voicemail-greeting",
		"Upload Busy Voicemail Greeting for Virtual Line",
		`Upload a busy voicemail greeting for a virtual line.

Examples:
  webex calling virtual-line-call update-busy-voicemail-greeting --virtual-line-id <id> --file greeting.wav`,
		"/telephony/config/virtualLines/{entityId}/voicemail/actions/uploadBusyGreeting/invoke",
		"virtual-line-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(virtualLineCallCmd,
		"update-no-answer-voicemail-greeting",
		"Upload No-Answer Voicemail Greeting for Virtual Line",
		`Upload a no-answer voicemail greeting for a virtual line.

Examples:
  webex calling virtual-line-call update-no-answer-voicemail-greeting --virtual-line-id <id> --file greeting.wav`,
		"/telephony/config/virtualLines/{entityId}/voicemail/actions/uploadNoAnswerGreeting/invoke",
		"virtual-line-id", func(cmd *cobra.Command) {})

	// workspace-call greetings
	registerGreetingUpload(workspaceCallCmd,
		"upload-intercept-announcement-file",
		"Upload Call Intercept Greeting for Workspace",
		`Upload a call intercept announcement greeting for a workspace.

Examples:
  webex calling workspace-call upload-intercept-announcement-file --workspace-id <id> --file greeting.wav`,
		"/telephony/config/workspaces/{entityId}/actions/announcementUpload/invoke",
		"workspace-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(workspaceCallCmd,
		"update-busy-voicemail-greeting-place",
		"Upload Busy Voicemail Greeting for Workspace",
		`Upload a busy voicemail greeting for a workspace (place).

Examples:
  webex calling workspace-call update-busy-voicemail-greeting-place --workspace-id <id> --file greeting.wav`,
		"/telephony/config/workspaces/{entityId}/voicemail/actions/uploadBusyGreeting/invoke",
		"workspace-id", func(cmd *cobra.Command) {})

	registerGreetingUpload(workspaceCallCmd,
		"update-no-answer-voicemail-greeting-place",
		"Upload No-Answer Voicemail Greeting for Workspace",
		`Upload a no-answer voicemail greeting for a workspace (place).

Examples:
  webex calling workspace-call update-no-answer-voicemail-greeting-place --workspace-id <id> --file greeting.wav`,
		"/telephony/config/workspaces/{entityId}/voicemail/actions/uploadNoAnswerGreeting/invoke",
		"workspace-id", func(cmd *cobra.Command) {})
}

// registerGreetingUpload creates a voicemail/intercept greeting upload command.
// entityFlag is the flag name for the entity ID (e.g. "person-id", "workspace-id").
// If empty, no entity ID flag is added (used for "me" endpoints).
// endpoint uses {entityId} as placeholder for the entity ID value.
func registerGreetingUpload(parentCmd *cobra.Command, use, short, long string,
	endpoint string, entityFlag string, extraFlags func(*cobra.Command)) {

	var filePath string
	var entityID string
	var orgID string

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

			// Check file size (API enforces 5000 KB max)
			info, err := f.Stat()
			if err != nil {
				return fmt.Errorf("reading file info: %w", err)
			}
			if info.Size() > 5000*1024 {
				return fmt.Errorf("file size %d bytes exceeds 5000 KB limit", info.Size())
			}

			// Build URL with entity ID substitution
			path := endpoint
			if entityFlag != "" {
				path = replaceEntityID(path, entityID)
			}
			url := config.CallingBaseURL + path
			if orgID != "" {
				url += "?orgId=" + orgID
			}

			parts := []audio.MultipartPart{
				{
					FieldName:   "file",
					FileName:    filepath.Base(filePath),
					ContentType: "audio/wav",
					Data:        f,
				},
			}

			body, statusCode, err := audio.UploadMultipart("POST", url, parts)
			if err != nil {
				return err
			}

			return output.Print(body, statusCode)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to WAV file to upload")
	cmd.MarkFlagRequired("file")
	if entityFlag != "" {
		cmd.Flags().StringVar(&entityID, entityFlag, "", entityFlag+" value")
		cmd.MarkFlagRequired(entityFlag)
	}
	cmd.Flags().StringVar(&orgID, "org-id", "", "Organization ID")
	extraFlags(cmd)

	parentCmd.AddCommand(cmd)
}

// replaceEntityID substitutes {entityId} in the path with the actual value.
func replaceEntityID(path, id string) string {
	result := path
	for _, placeholder := range []string{"{entityId}"} {
		if idx := len(result); idx > 0 {
			result = replaceFirst(result, placeholder, id)
		}
	}
	return result
}

func replaceFirst(s, old, new string) string {
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}
