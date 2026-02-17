package device

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

var deviceConfigurationsCmd = &cobra.Command{
	Use:   "device-configurations",
	Short: "DeviceConfigurations commands",
}

func init() {
	cmd.DeviceCmd.AddCommand(deviceConfigurationsCmd)

	{ // list
		var deviceId string
		var key string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Device Configurations for device",
			Long:  `Lists all device configurations associated with the given device ID. Administrators can list configurations for all devices within an organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/deviceConfigurations")
				req.QueryParam("deviceId", deviceId)
				req.QueryParam("key", key)
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
		cmd.Flags().StringVar(&deviceId, "device-id", "", "List device configurations by device ID.")
		cmd.Flags().StringVar(&key, "key", "", "This can optionally be used to filter configurations. Keys are composed of segments. It's possible to use absolute paths, wildcards or ranges.  - **Absolute** gives only one configuration as a result. `Conference.MaxReceiveCallRate` for example gives the Conference `MaxReceiveCallRate` configuration.  + **Wildcards** (\\*) can specify multiple configurations with shared segments. `Audio.Ultrasound.*` for example will filter on all Audio Ultrasound configurations.  - **Range** ([_number_]) can be used to filter numbered segments. `FacilityService.Service[1].Name` for instance only shows the first `FacilityService` Service Name configuration, `FacilityService.Service[*].Name` shows all, `FacilityService.Service[1..3].Name` shows the first three and `FacilityService.Service[2..n].Name` shows all starting at 2. Note that [RFC 3986 3.2.2](https://www.ietf.org/rfc/rfc3986.html#section-3.2.2) does not allow square brackets in urls outside the host, so to specify range in a configuration key you will need to encode them to %5B for [ and %5D for ].")
		deviceConfigurationsCmd.AddCommand(cmd)
	}

	{ // update
		var deviceId string
		var op string
		var path string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update Device Configurations",
			Long:  `Edit configurations for the device specified by device ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/deviceConfigurations")
				req.QueryParam("deviceId", deviceId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("op", op)
					req.BodyString("path", path)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&deviceId, "device-id", "", "Update device configurations by device ID.")
		cmd.Flags().StringVar(&op, "op", "", "")
		cmd.Flags().StringVar(&path, "path", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		deviceConfigurationsCmd.AddCommand(cmd)
	}

}
