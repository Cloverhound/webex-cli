package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Scope preset definitions.
var (
	userScopes = []string{
		"spark:all",
		"spark:calls_read",
		"spark:calls_write",
		"spark:memberships_read",
		"spark:memberships_write",
		"spark:messages_read",
		"spark:messages_write",
		"spark:rooms_read",
		"spark:rooms_write",
		"spark:teams_read",
		"spark:teams_write",
		"spark:team_memberships_read",
		"spark:team_memberships_write",
		"spark:people_read",
		"spark:people_write",
		"spark:devices_read",
		"spark:devices_write",
		"spark:places_read",
		"spark:places_write",
		"spark:organizations_read",
		"spark:recordings_read",
		"spark:recordings_write",
		"spark:telephony_config_read",
		"spark:telephony_config_write",
		"spark:xsi",
		"spark:webrtc_calling",
		"meeting:recordings_read",
		"meeting:recordings_write",
		"meeting:preferences_read",
		"meeting:preferences_write",
		"meeting:schedules_read",
		"meeting:schedules_write",
		"meeting:participants_read",
		"meeting:participants_write",
		"meeting:controls_read",
		"meeting:controls_write",
		"meeting:transcripts_read",
		"meeting:summaries_read",
		"meeting:summaries_write",
	}

	adminReadScopes = []string{
		"spark-admin:people_read",
		"spark-admin:licenses_read",
		"spark-admin:roles_read",
		"spark-admin:organizations_read",
		"spark-admin:devices_read",
		"spark-admin:places_read",
		"spark-admin:locations_read",
		"spark-admin:workspaces_read",
		"spark-admin:workspace_locations_read",
		"spark-admin:workspace_metrics_read",
		"spark-admin:telephony_config_read",
		"spark-admin:telephony_pstn_read",
		"spark-admin:recordings_read",
		"spark-admin:reports_read",
		"spark-admin:calling_cdr_read",
		"spark-admin:call_qualities_read",
		"spark-admin:hybrid_clusters_read",
		"spark-admin:hybrid_connectors_read",
		"spark-admin:resource_groups_read",
		"spark-admin:resource_group_memberships_read",
		"spark-admin:broadworks_enterprises_read",
		"spark-admin:broadworks_subscribers_read",
		"spark-admin:broadworks_billing_reports_read",
		"spark-admin:wholesale_customers_read",
		"spark-admin:wholesale_subscribers_read",
		"spark-admin:wholesale_billing_reports_read",
		"spark-admin:wholesale_sub_partners_read",
		"spark-admin:messages_read",
		"meeting:admin_recordings_read",
		"meeting:admin_transcripts_read",
		"meeting:admin_participants_read",
		"meeting:admin_schedule_read",
		"meeting:admin_preferences_read",
		"meeting:admin_config_read",
		"analytics:read_all",
		"audit:events_read",
		"spark-compliance:events_read",
		"spark-compliance:memberships_read",
		"spark-compliance:messages_read",
		"spark-compliance:rooms_read",
		"spark-compliance:teams_read",
		"spark-compliance:team_memberships_read",
		"spark-compliance:recordings_read",
		"spark-compliance:meetings_read",
		"spark-compliance:webhooks_read",
		"identity:people_read",
		"identity:groups_read",
		"identity:organizations_read",
		"identity:tokens_read",
		"cjp:config_read",
		"cjds:admin_org_read",
	}

	adminWriteScopes = []string{
		"spark-admin:people_write",
		"spark-admin:organizations_write",
		"spark-admin:devices_write",
		"spark-admin:places_write",
		"spark-admin:locations_write",
		"spark-admin:workspaces_write",
		"spark-admin:workspace_locations_write",
		"spark-admin:telephony_config_write",
		"spark-admin:telephony_pstn_write",
		"spark-admin:recordings_write",
		"spark-admin:reports_write",
		"spark-admin:calls_write",
		"spark-admin:resource_group_memberships_write",
		"spark-admin:broadworks_enterprises_write",
		"spark-admin:broadworks_subscribers_write",
		"spark-admin:broadworks_billing_reports_write",
		"spark-admin:wholesale_customers_write",
		"spark-admin:wholesale_subscribers_write",
		"spark-admin:wholesale_billing_reports_write",
		"spark-admin:wholesale_sub_partners_write",
		"spark-admin:wholesale_workspace_write",
		"spark-admin:messages_write",
		"meeting:admin_recordings_write",
		"meeting:admin_schedule_write",
		"meeting:admin_preferences_write",
		"meeting:admin_config_write",
		"spark-compliance:memberships_write",
		"spark-compliance:messages_write",
		"spark-compliance:rooms_write",
		"spark-compliance:team_memberships_write",
		"spark-compliance:recordings_write",
		"spark-compliance:meetings_write",
		"spark-compliance:webhooks_write",
		"identity:people_rw",
		"identity:groups_rw",
		"identity:organizations_rw",
		"Identity:contact",
		"identity:tokens_write",
		"identity:placeonetimepassword_create",
		"Identity:one_time_password",
		"cjp:config",
		"cjp:config_write",
		"cjp:user",
		"cjds:admin_org_write",
		"cloud-contact-center:pod_conv",
	}

	// Scopes to exclude from the final list.
	excludedScopes = map[string]bool{
		"spark:kms":       true,
		"spark-admin:xsi": true,
	}
)

const defaultRedirectURI = "http://localhost:8085/callback"

var configCreateIntegrationCmd = &cobra.Command{
	Use:   "create-integration",
	Short: "Create a Webex Integration and configure the CLI",
	Long: `Create a new Webex Integration via the Webex API and automatically configure
the CLI with the resulting client-id, client-secret, and scopes.

Requires a personal access token from https://developer.webex.com (My Webex Apps).
The token is used only during this command to create the integration and is not stored.

After creation, run 'webex login' to authenticate via OAuth with the new integration.`,
	RunE: runCreateIntegration,
}

func init() {
	f := configCreateIntegrationCmd.Flags()
	f.String("token", "", "Personal access token (skips interactive prompt)")
	f.String("name", "", "Integration name")
	f.String("description", "", "Integration description")
	f.String("contact-email", "", "Contact email")
	f.String("company-name", "", "Company name")
	f.String("redirect-uris", "", "Comma-separated redirect URIs (default: http://localhost:8085/callback)")
	f.String("scope-preset", "", "Scope preset: user, admin-read, admin-readwrite")
	f.String("scopes", "", "Space-separated scopes (manual, overrides --scope-preset)")
	f.Bool("non-interactive", false, "Run without prompts (requires all flags)")
}

func runCreateIntegration(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)
	nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

	// 1. Get access token
	token, err := getToken(cmd, reader, nonInteractive)
	if err != nil {
		return err
	}

	// 2. Fetch org info and confirm
	orgID, orgName, err := fetchAndConfirmOrg(token, reader, nonInteractive)
	if err != nil {
		return err
	}

	// 3. Collect integration details
	name, err := getRequiredInput(cmd, reader, "name", "Integration name: ", nonInteractive)
	if err != nil {
		return err
	}
	description, err := getRequiredInput(cmd, reader, "description", "Description: ", nonInteractive)
	if err != nil {
		return err
	}
	contactEmail, err := getRequiredInput(cmd, reader, "contact-email", "Contact email: ", nonInteractive)
	if err != nil {
		return err
	}
	companyName, err := getRequiredInput(cmd, reader, "company-name", "Company name: ", nonInteractive)
	if err != nil {
		return err
	}

	// 4. Redirect URIs
	redirectURIs, err := getRedirectURIs(cmd, reader, nonInteractive)
	if err != nil {
		return err
	}

	// 5. Scopes
	selectedScopes, err := getScopes(cmd, reader, nonInteractive)
	if err != nil {
		return err
	}

	// 6. Create integration via API
	clientID, clientSecret, integrationID, err := createIntegration(
		token, orgID, name, description, contactEmail, companyName,
		redirectURIs, selectedScopes,
	)
	if err != nil {
		return err
	}

	// 7. Display results and auto-configure
	fmt.Println()
	fmt.Println("Integration created successfully!")
	fmt.Printf("  ID:            %s\n", integrationID)
	fmt.Printf("  Client ID:     %s\n", clientID)
	fmt.Printf("  Client Secret: %s\n", clientSecret)
	fmt.Println()

	cfg, err := appconfig.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	cfg.ClientID = clientID
	cfg.ClientSecret = clientSecret
	cfg.Scopes = strings.Join(selectedScopes, " ")
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Println("Configuring CLI...")
	fmt.Println("  Set client-id")
	fmt.Println("  Set client-secret")
	fmt.Println("  Set scopes")
	fmt.Println()
	fmt.Printf("Organization: %s\n", orgName)
	fmt.Println()
	fmt.Println("Ready! Run 'webex login' to authenticate with your new integration.")

	return nil
}

func getToken(cmd *cobra.Command, reader *bufio.Reader, nonInteractive bool) (string, error) {
	flagToken, _ := cmd.Flags().GetString("token")
	if flagToken != "" {
		return flagToken, nil
	}
	if nonInteractive {
		return "", fmt.Errorf("--token is required in non-interactive mode")
	}

	fmt.Print("Enter your Webex access token (from developer.webex.com): ")
	tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after masked input
	if err != nil {
		return "", fmt.Errorf("reading token: %w", err)
	}
	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}
	return token, nil
}

type personResponse struct {
	OrgID string `json:"orgId"`
}

type orgResponse struct {
	DisplayName string `json:"displayName"`
}

func fetchAndConfirmOrg(token string, reader *bufio.Reader, nonInteractive bool) (string, string, error) {
	// GET /v1/people/me
	req, err := http.NewRequest("GET", "https://webexapis.com/v1/people/me", nil)
	if err != nil {
		return "", "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fetching user info: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch user info (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var person personResponse
	if err := json.Unmarshal(body, &person); err != nil {
		return "", "", fmt.Errorf("parsing user info: %w", err)
	}

	// GET /v1/organizations/{orgId}
	req, err = http.NewRequest("GET", "https://webexapis.com/v1/organizations/"+person.OrgID, nil)
	if err != nil {
		return "", "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fetching org info: %w", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch org info (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var org orgResponse
	if err := json.Unmarshal(body, &org); err != nil {
		return "", "", fmt.Errorf("parsing org info: %w", err)
	}

	fmt.Printf("Organization: %s (%s)\n", org.DisplayName, person.OrgID)

	if !nonInteractive {
		fmt.Print("Proceed with this organization? (y/N): ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			return "", "", fmt.Errorf("cancelled by user")
		}
	}

	fmt.Println()
	return person.OrgID, org.DisplayName, nil
}

func getRequiredInput(cmd *cobra.Command, reader *bufio.Reader, flagName, prompt string, nonInteractive bool) (string, error) {
	flagVal, _ := cmd.Flags().GetString(flagName)
	if flagVal != "" {
		return flagVal, nil
	}
	if nonInteractive {
		return "", fmt.Errorf("--%s is required in non-interactive mode", flagName)
	}

	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("%s cannot be empty", flagName)
	}
	return input, nil
}

func getRedirectURIs(cmd *cobra.Command, reader *bufio.Reader, nonInteractive bool) ([]string, error) {
	flagVal, _ := cmd.Flags().GetString("redirect-uris")
	if flagVal != "" {
		uris := parseCSV(flagVal)
		// Ensure default is included
		if !contains(uris, defaultRedirectURI) {
			uris = append([]string{defaultRedirectURI}, uris...)
		}
		return uris, nil
	}
	if nonInteractive {
		return []string{defaultRedirectURI}, nil
	}

	fmt.Printf("Redirect URIs (comma-separated, default: %s): ", defaultRedirectURI)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{defaultRedirectURI}, nil
	}

	uris := parseCSV(input)
	if !contains(uris, defaultRedirectURI) {
		uris = append([]string{defaultRedirectURI}, uris...)
	}
	return uris, nil
}

func getScopes(cmd *cobra.Command, reader *bufio.Reader, nonInteractive bool) ([]string, error) {
	// Manual scopes flag takes priority
	manualScopes, _ := cmd.Flags().GetString("scopes")
	if manualScopes != "" {
		return filterScopes(strings.Fields(manualScopes)), nil
	}

	// Scope preset flag
	preset, _ := cmd.Flags().GetString("scope-preset")
	if preset != "" {
		return resolveScopePreset(preset)
	}

	if nonInteractive {
		return nil, fmt.Errorf("--scope-preset or --scopes is required in non-interactive mode")
	}

	// Interactive scope selection
	fmt.Println()
	fmt.Println("Available scope presets:")
	fmt.Println("  1) User — Personal access scopes")
	fmt.Println("  2) Admin Read-Only — Organization admin read access")
	fmt.Println("  3) Admin Read/Write — Full organization admin access")
	fmt.Println("  4) Enter scopes manually")
	fmt.Println()

	for {
		fmt.Print("Select a preset (or 'check <number>' to preview): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.HasPrefix(input, "check ") {
			numStr := strings.TrimPrefix(input, "check ")
			num, err := strconv.Atoi(strings.TrimSpace(numStr))
			if err != nil || num < 1 || num > 3 {
				fmt.Println("Invalid option. Use 'check 1', 'check 2', or 'check 3'.")
				continue
			}
			previewScopes(num)
			fmt.Println()
			continue
		}

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > 4 {
			fmt.Println("Invalid selection. Enter 1-4.")
			continue
		}

		if num == 4 {
			fmt.Print("Enter scopes (space-separated): ")
			scopeInput, _ := reader.ReadString('\n')
			scopeInput = strings.TrimSpace(scopeInput)
			if scopeInput == "" {
				fmt.Println("Scopes cannot be empty.")
				continue
			}
			return filterScopes(strings.Fields(scopeInput)), nil
		}

		presetNames := []string{"", "user", "admin-read", "admin-readwrite"}
		return resolveScopePreset(presetNames[num])
	}
}

func resolveScopePreset(preset string) ([]string, error) {
	var scopes []string
	switch preset {
	case "user":
		scopes = append(scopes, userScopes...)
	case "admin-read":
		scopes = append(scopes, userScopes...)
		scopes = append(scopes, adminReadScopes...)
	case "admin-readwrite":
		scopes = append(scopes, userScopes...)
		scopes = append(scopes, adminReadScopes...)
		scopes = append(scopes, adminWriteScopes...)
	default:
		return nil, fmt.Errorf("unknown scope preset: %s (valid: user, admin-read, admin-readwrite)", preset)
	}
	return filterScopes(scopes), nil
}

func previewScopes(num int) {
	var scopes []string
	var label string
	switch num {
	case 1:
		label = "User"
		scopes = userScopes
	case 2:
		label = "Admin Read-Only"
		scopes = append(append([]string{}, userScopes...), adminReadScopes...)
	case 3:
		label = "Admin Read/Write"
		scopes = append(append(append([]string{}, userScopes...), adminReadScopes...), adminWriteScopes...)
	}
	fmt.Printf("\n%s scopes (%d):\n", label, len(filterScopes(scopes)))
	for _, s := range filterScopes(scopes) {
		fmt.Printf("  %s\n", s)
	}
}

func filterScopes(scopes []string) []string {
	var filtered []string
	for _, s := range scopes {
		if !excludedScopes[s] {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

type createIntegrationRequest struct {
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ContactEmail string   `json:"contactEmail"`
	CompanyName  string   `json:"companyName"`
	RedirectURLs []string `json:"redirectUrls"`
	Scopes       []string `json:"scopes"`
}

type createIntegrationResponse struct {
	ID           string `json:"id"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func createIntegration(token, orgID, name, description, contactEmail, companyName string, redirectURIs, scopes []string) (clientID, clientSecret, integrationID string, err error) {
	payload := createIntegrationRequest{
		Type:         "integration",
		Name:         name,
		Description:  description,
		ContactEmail: contactEmail,
		CompanyName:  companyName,
		RedirectURLs: redirectURIs,
		Scopes:       scopes,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", "", "", fmt.Errorf("marshaling request: %w", err)
	}

	url := fmt.Sprintf("https://webexapis.com/v1/applications?orgId=%s", orgID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", "", "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("creating integration: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", "", "", fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result createIntegrationResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", "", "", fmt.Errorf("parsing response: %w", err)
	}

	if result.ClientID == "" || result.ClientSecret == "" {
		return "", "", "", fmt.Errorf("unexpected response: missing clientId or clientSecret\n%s", string(respBody))
	}

	return result.ClientID, result.ClientSecret, result.ID, nil
}

func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func contains(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}
