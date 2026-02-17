package calling

import (
	"fmt"
	"strings"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = config.Token
var _ = output.Print
var _ = strings.Join

var conferenceControlsCmd = &cobra.Command{
	Use:   "conference-controls",
	Short: "ConferenceControls commands",
}

func init() {
	cmd.CallingCmd.AddCommand(conferenceControlsCmd)

	{ // start
		var callIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "start",
			Short: "Start Conference",
			Long: `Join the user's calls into a conference.  A minimum of two call IDs are required. Each call ID identifies an existing call between the user
and a participant to be added to the conference.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("callIds", callIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringSliceVar(&callIds, "call-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // release
		cmd := &cobra.Command{
			Use:   "release",
			Short: "Release Conference",
			Long:  `Release the conference (the host and all participants). Note that for a 3WC (three-way call) the [Transfer API](/docs/api/v1/call-controls/transfer) can be used to perform an attended transfer so that the participants remain connected.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/conference")
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // get
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Conference Details",
			Long:  `Get the details of the conference.  An empty JSON object body is returned if there is no conference.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/conference")
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(true)
					if err != nil {
						return err
					}
					return output.Print(resp, statusCode)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // create-participant
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-participant",
			Short: "Add Participant",
			Long:  `Adds a participant to an existing conference.  The request body contains the participant's call ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/addParticipant")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // mute
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "mute",
			Short: "Mute",
			Long:  `Mutes the host or a participant. Mutes the host when no request body is provided (i.e. media stream from the host will not be transmitted to the conference).  Mutes a participant when the request body contains the participant's call ID (i.e. media stream from the participant will not be transmitted to the conference).`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/mute")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // unmute
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "unmute",
			Short: "Unmute",
			Long:  `Unmutes the host or a participant. Unmutes the host when no request body is provided (i.e. media stream from the host will be transmitted to the conference).  Unmutes a participant when the request body contains the participant's call ID (i.e. media stream from the participant will be transmitted to the conference).`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/unmute")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // deafen-participant
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "deafen-participant",
			Short: "Deafen Participant",
			Long: `Deafens a participant (i.e. media stream will not be transmitted to the participant).
The request body contains the call ID of the participant to deafen`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/deafen")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // undeafen-participant
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "undeafen-participant",
			Short: "Undeafen Participant",
			Long: `Undeafens a participant (i.e. resume transmitting the conference media stream to the participant). 
The request body contains the call ID of the participant to undeafen.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/undeafen")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // hold
		cmd := &cobra.Command{
			Use:   "hold",
			Short: "Hold",
			Long:  `Hold the conference host.  There is no request body.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/hold")
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		conferenceControlsCmd.AddCommand(cmd)
	}

	{ // resume
		cmd := &cobra.Command{
			Use:   "resume",
			Short: "Resume",
			Long:  `Resumes the held conference host.  There is no request body.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/conference/resume")
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		conferenceControlsCmd.AddCommand(cmd)
	}

}
