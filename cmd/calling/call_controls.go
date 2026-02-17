package calling

import (
	"fmt"

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

var callControlsCmd = &cobra.Command{
	Use:   "call-controls",
	Short: "CallControls commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callControlsCmd)

	{ // dial
		var destination string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "dial",
			Short: "Dial",
			Long:  `Initiate an outbound call to a specified destination. This is also commonly referred to as Click to Call or Click to Dial. Alerts occur on all the devices belonging to a user unless an optional endpointId is specified in which case only the device or application identified by the endpointId is alerted. When a user answers an alerting device, an outbound call is placed from that device to the destination.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/dial")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("destination", destination)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // answer
		var callId string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "answer",
			Short: "Answer",
			Long:  `Answer an incoming call. When no endpointId is specified, the call is answered on the user's primary device. When an endpointId is specified, the call is answered on the device or application identified by the endpointId. The answer API is rejected if the device is not alerting for the call or the device does not support answer via API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/answer")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // reject
		var callId string
		var action string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "reject",
			Short: "Reject",
			Long:  `Reject an unanswered incoming call.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/reject")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("action", action)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // hangup
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "hangup",
			Short: "Hangup",
			Long:  `Hangup a call. If used on an unanswered incoming call, the call is rejected and sent to busy.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/hangup")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // hold
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "hold",
			Short: "Hold",
			Long:  `Hold a connected call.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/hold")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // resume
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "resume",
			Short: "Resume",
			Long:  `Resume a held call.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/resume")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // mute
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "mute",
			Short: "Mute",
			Long:  `Mute a call. This API can only be used for a call that reports itself as mute capable via the muteCapable field in the call details.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/mute")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // unmute
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "unmute",
			Short: "Unmute",
			Long:  `Unmute a call. This API can only be used for a call that reports itself as mute capable via the muteCapable field in the call details.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/unmute")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // divert
		var callId string
		var destination string
		var toVoicemail bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "divert",
			Short: "Divert",
			Long:  `Divert a call to a destination or a user's voicemail. This is also commonly referred to as a Blind Transfer.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/divert")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("destination", destination)
					req.BodyBool("toVoicemail", toVoicemail, cmd.Flags().Changed("to-voicemail"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().BoolVar(&toVoicemail, "to-voicemail", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // transfer
		var callId1 string
		var callId2 string
		var destination string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "transfer",
			Short: "Transfer",
			Long:  "Transfer two calls together.\n\nUnanswered incoming calls cannot be transferred but can be diverted using the divert API.\n\nIf the user has only two calls and wants to transfer them together, the `callId1` and `callId2` parameters are optional and when not provided the calls are automatically selected and transferred.\n\nIf the user has more than two calls and wants to transfer two of them together, the `callId1` and `callId2` parameters are mandatory to specify which calls are being transferred. Those are also commonly referred to as Attended Transfer, Consultative Transfer, or Supervised Transfer and will return a `204` response.\n\nIf the user wants to transfer one call to a new destination but only when the destination responds, the `callId1` and destination parameters are mandatory to specify the call being transferred and the destination.\n\nThis is referred to as a Mute Transfer and is similar to the divert API with the difference of waiting for the destination to respond prior to transferring the call. If the destination does not respond, the call is not transferred. This will return a `201` response.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/transfer")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId1", callId1)
					req.BodyString("callId2", callId2)
					req.BodyString("destination", destination)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId1, "call-id1", "", "")
		cmd.Flags().StringVar(&callId2, "call-id2", "", "")
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // park
		var callId string
		var destination string
		var isGroupPark bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "park",
			Short: "Park",
			Long:  `Park a connected call. The number field in the response can be used as the destination for the retrieve command to retrieve the parked call.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/park")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("destination", destination)
					req.BodyBool("isGroupPark", isGroupPark, cmd.Flags().Changed("is-group-park"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().BoolVar(&isGroupPark, "is-group-park", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // get
		var destination string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Retrieve",
			Long:  `Retrieve a parked call. A new call is initiated to perform the retrieval in a similar manner to the dial command. The number field from the park command response can be used as the destination for the retrieve command.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/retrieve")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("destination", destination)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // start-recording
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "start-recording",
			Short: "Start Recording",
			Long:  `Start recording a call. Use of this API is only valid when the user's call recording mode is set to "On Demand".`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/startRecording")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // stop-recording
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "stop-recording",
			Short: "Stop Recording",
			Long:  `Stop recording a call. Use of this API is only valid when a call is being recorded and the user's call recording mode is set to "On Demand".`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/stopRecording")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // pause-recording
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "pause-recording",
			Short: "Pause Recording",
			Long:  `Pause recording on a call. Use of this API is only valid when a call is being recorded and the user's call recording mode is set to "On Demand" or "Always with Pause/Resume".`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/pauseRecording")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // resume-recording
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "resume-recording",
			Short: "Resume Recording",
			Long:  `Resume recording a call. Use of this API is only valid when a call's recording is paused and the user's call recording mode is set to "On Demand" or "Always with Pause/Resume".`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/resumeRecording")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // transmit-dtmf
		var callId string
		var dtmf string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "transmit-dtmf",
			Short: "Transmit DTMF",
			Long:  `Transmit DTMF digits to a call.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/transmitDtmf")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("dtmf", dtmf)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&dtmf, "dtmf", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // push
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "push",
			Short: "Push",
			Long:  `Pushes a call from the assistant to the executive the call is associated with. Use of this API is only valid when the assistant's call is associated with an executive.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/push")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // pickup
		var target string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "pickup",
			Short: "Pickup",
			Long:  `Picks up an incoming call to another user. A new call is initiated to perform the pickup in a similar manner to the dial command. When target is not present, the API pickups up a call from the user's call pickup group. When target is present, the API pickups an incoming call from the specified target user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/pickup")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("target", target)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&target, "target", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // barge
		var target string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "barge",
			Short: "Barge In",
			Long:  `Barge-in on another user's answered call. A new call is initiated to perform the barge-in in a similar manner to the dial command.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/bargeIn")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("target", target)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&target, "target", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // list
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Calls",
			Long:  `Get the list of details for all active calls associated with the user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/calls")
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
		callControlsCmd.AddCommand(cmd)
	}

	{ // get-2
		var callId string
		cmd := &cobra.Command{
			Use:   "get-2",
			Short: "Get Call Details",
			Long:  `Get the details of the specified active call for the user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/calls/{callId}")
				req.PathParam("callId", callId)
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
		cmd.Flags().StringVar(&callId, "call-id", "", "The call identifier of the call.")
		cmd.MarkFlagRequired("call-id")
		callControlsCmd.AddCommand(cmd)
	}

	{ // list-history
		var typeVal string
		cmd := &cobra.Command{
			Use:   "list-history",
			Short: "List Call History",
			Long:  "Get the list of call history records for the user. A maximum of 20 call history records per type (`placed`, `missed`, `received`) are returned.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/calls/history")
				req.QueryParam("type", typeVal)
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
		cmd.Flags().StringVar(&typeVal, "type", "", "The type of call history records to retrieve. If not specified, then all call history records are retrieved.")
		callControlsCmd.AddCommand(cmd)
	}

	{ // dial-member-id
		var memberId string
		var orgId string
		var destination string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "dial-member-id",
			Short: "Dial by Member ID",
			Long:  `Initiate an outbound call to a specified destination. This is also commonly referred to as Click to Call or Click to Dial. Alerts occur on all the devices belonging to a user unless an optional endpointId is specified in which case only the device or application identified by the endpointId is alerted. When a user answers an alerting device, an outbound call is placed from that device to the destination.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/members/{memberId}/dial")
				req.PathParam("memberId", memberId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("destination", destination)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&memberId, "member-id", "", "Unique identifier for the member. Member ID can be one of the following types: PEOPLE, PLACE, or VIRTUAL_LINE")
		cmd.MarkFlagRequired("member-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Id of the organization to which the member belongs. If not provided, the orgId of the Service App is used. If provided, the organization must be the same as or managed by the Service App's organization.")
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // answer-member-id
		var memberId string
		var orgId string
		var callId string
		var endpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "answer-member-id",
			Short: "Answer by Member ID",
			Long:  `Answer an incoming call. When no endpointId is specified, the call is answered on the user's primary device. When an endpointId is specified, the call is answered on the device or application identified by the endpointId. The answer API is rejected if the device is not alerting for the call or the device does not support answer via API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/members/{memberId}/answer")
				req.PathParam("memberId", memberId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callId", callId)
					req.BodyString("endpointId", endpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&memberId, "member-id", "", "Unique identifier for the member. Member ID can be one of the following types: PEOPLE, PLACE, or VIRTUAL_LINE")
		cmd.MarkFlagRequired("member-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Id of the organization to which the member belongs. If not provided, the orgId of the Service App is used. If provided, the organization must be the same as or managed by the Service App's organization.")
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&endpointId, "endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // hangup-member-id
		var memberId string
		var orgId string
		var callId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "hangup-member-id",
			Short: "Hangup by Member ID",
			Long:  `Hangup a call. If used on an unanswered incoming call, the call is rejected and sent to busy.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/calls/members/{memberId}/hangup")
				req.PathParam("memberId", memberId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&memberId, "member-id", "", "Unique identifier for the member. Member ID can be one of the following types: PEOPLE, PLACE, or VIRTUAL_LINE")
		cmd.MarkFlagRequired("member-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Id of the organization to which the member belongs. If not provided, the orgId of the Service App is used. If provided, the organization must be the same as or managed by the Service App's organization.")
		cmd.Flags().StringVar(&callId, "call-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callControlsCmd.AddCommand(cmd)
	}

	{ // list-member-id
		var memberId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-member-id",
			Short: "List Calls by Member ID",
			Long:  `Get the list of details for all active calls associated with the member.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/calls/members/{memberId}/calls")
				req.PathParam("memberId", memberId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&memberId, "member-id", "", "Unique identifier for the member. Member ID can be one of the following types: PEOPLE, PLACE, or VIRTUAL_LINE")
		cmd.MarkFlagRequired("member-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Id of the organization to which the member belongs. If not provided, the orgId of the Service App is used. If provided, the organization must be the same as or managed by the Service App's organization.")
		callControlsCmd.AddCommand(cmd)
	}

	{ // get-member-id
		var memberId string
		var callId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-member-id",
			Short: "Get Call Details by Member ID",
			Long:  `Get the details of the specified active call for the member.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/calls/members/{memberId}/calls/{callId}")
				req.PathParam("memberId", memberId)
				req.PathParam("callId", callId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&memberId, "member-id", "", "Unique identifier for the member. Member ID can be one of the following types: PEOPLE, PLACE, or VIRTUAL_LINE")
		cmd.MarkFlagRequired("member-id")
		cmd.Flags().StringVar(&callId, "call-id", "", "The call identifier of the call.")
		cmd.MarkFlagRequired("call-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Id of the organization to which the member belongs. If not provided, the orgId of the Service App is used. If provided, the organization must be the same as or managed by the Service App's organization.")
		callControlsCmd.AddCommand(cmd)
	}

}
