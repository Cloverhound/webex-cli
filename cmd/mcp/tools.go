package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

// msgSearchSpaceLimit is the default number of recently-active spaces to scan
// during a keyword message search. Webex has no global full-text search API,
// so we fan out across the N most-active rooms and filter client-side.
const msgSearchSpaceLimit = 15

// registerTools adds all 10 Webex MCP tools to the server.
func registerTools(s *server.MCPServer) {
	s.AddTool(
		mcplib.NewTool("webex_list_spaces",
			mcplib.WithDescription(
				"List Webex spaces (rooms) the authenticated user is a member of. "+
					"Optionally filter by name keyword and/or space type ('direct' or 'group'). "+
					"Use this to discover room_id values before calling webex_get_messages or webex_send_message.",
			),
			mcplib.WithString("keyword",
				mcplib.Description("Case-insensitive substring to filter by space name (e.g. 'Acme', 'Q2 Deal')"),
			),
			mcplib.WithString("space_type",
				mcplib.Description("Filter by type: 'direct' (1-on-1 DMs) or 'group' (group spaces). Omit for all."),
				mcplib.Enum("direct", "group"),
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Maximum number of spaces to return (1–200, default 20)"),
			),
		),
		handleListSpaces,
	)

	s.AddTool(
		mcplib.NewTool("webex_get_messages",
			mcplib.WithDescription(
				"Retrieve recent messages from a Webex space by room ID or partial room name. "+
					"Messages are returned newest-first. "+
					"When precision matters, get an exact room_id via webex_list_spaces first.",
			),
			mcplib.WithString("room_id",
				mcplib.Description("Exact Webex room/space ID (preferred). Provide this OR room_name."),
			),
			mcplib.WithString("room_name",
				mcplib.Description("Partial space name — the first match is used. Provide this OR room_id."),
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Number of recent messages to return (1–100, default 20)"),
			),
		),
		handleGetMessages,
	)

	s.AddTool(
		mcplib.NewTool("webex_search_messages",
			mcplib.WithDescription(
				"Search for a keyword or phrase across recently-active Webex spaces. "+
					"Scans the N most-recently-active spaces concurrently and filters client-side "+
					"(the Webex API has no global full-text search endpoint). "+
					"Increase space_limit for broader coverage at the cost of latency.",
			),
			mcplib.WithString("keyword",
				mcplib.Required(),
				mcplib.Description("Word or phrase to search for (min 2 chars)"),
			),
			mcplib.WithNumber("space_limit",
				mcplib.Description(fmt.Sprintf(
					"Number of recently-active spaces to scan (1–30, default %d). "+
						"Increase if results seem incomplete.", msgSearchSpaceLimit)),
			),
			mcplib.WithNumber("messages_per_space",
				mcplib.Description("Recent messages to inspect per space before keyword filtering (10–100, default 50)"),
			),
		),
		handleSearchMessages,
	)

	s.AddTool(
		mcplib.NewTool("webex_send_message",
			mcplib.WithDescription(
				"Send a plain-text (or Markdown) message to a Webex space or directly to a person by email. "+
					"Provide EITHER room_id (post to an existing space) OR to_person_email (start/continue a DM). "+
					"This is the only write-capable tool in this server.",
			),
			mcplib.WithString("text",
				mcplib.Required(),
				mcplib.Description("Plain-text message body (required, max 7 439 chars)"),
			),
			mcplib.WithString("room_id",
				mcplib.Description("Target space/room ID. Provide this OR to_person_email."),
			),
			mcplib.WithString("to_person_email",
				mcplib.Description("Recipient email for a direct message. Provide this OR room_id."),
			),
			mcplib.WithString("markdown",
				mcplib.Description("Optional Markdown version shown by clients that support rich formatting"),
			),
		),
		handleSendMessage,
	)

	s.AddTool(
		mcplib.NewTool("webex_list_direct_messages",
			mcplib.WithDescription(
				"List recent 1-on-1 (direct) Webex conversations with a last-message preview. "+
					"Ideal for a daily digest of incoming client messages.",
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Number of recent DM conversations to return (1–50, default 10)"),
			),
		),
		handleListDirectMessages,
	)

	s.AddTool(
		mcplib.NewTool("webex_get_person",
			mcplib.WithDescription(
				"Look up a Webex user by email address or display name. "+
					"Email provides an exact single match; display name performs a partial search "+
					"and may return several candidates.",
			),
			mcplib.WithString("email",
				mcplib.Description("Exact email address (most precise). Provide this OR display_name."),
			),
			mcplib.WithString("display_name",
				mcplib.Description("Display name to search (partial match, may return multiple). Provide this OR email."),
			),
		),
		handleGetPerson,
	)

	s.AddTool(
		mcplib.NewTool("webex_list_meetings",
			mcplib.WithDescription(
				"List upcoming and recent Webex meetings within an optional date range. "+
					"Defaults to 7 days ago through 30 days ahead. "+
					"Returned meeting IDs feed into webex_get_meeting_transcript and webex_get_meeting_summary.",
			),
			mcplib.WithString("from_date",
				mcplib.Description("Start of date range in ISO 8601 (e.g. '2025-04-01T00:00:00Z'). Defaults to 7 days ago."),
			),
			mcplib.WithString("to_date",
				mcplib.Description("End of date range in ISO 8601 (e.g. '2025-04-30T23:59:59Z'). Defaults to 30 days ahead."),
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Maximum meetings to return (1–100, default 20)"),
			),
		),
		handleListMeetings,
	)

	s.AddTool(
		mcplib.NewTool("webex_get_meeting_transcript",
			mcplib.WithDescription(
				"Retrieve the transcript for a completed Webex meeting. "+
					"The meeting must have ended with transcription enabled. "+
					"Use webex_list_meetings to find valid meeting IDs.",
			),
			mcplib.WithString("meeting_id",
				mcplib.Required(),
				mcplib.Description("Webex meeting ID — from webex_list_meetings or a meeting URL"),
			),
		),
		handleGetMeetingTranscript,
	)

	s.AddTool(
		mcplib.NewTool("webex_get_meeting_summary",
			mcplib.WithDescription(
				"Retrieve the Webex AI-generated summary (key points, action items, decisions) "+
					"for a completed meeting. "+
					"The meeting must be fully processed (status 'ended') before a summary is available. "+
					"For verbatim detail, use webex_get_meeting_transcript instead.",
			),
			mcplib.WithString("meeting_id",
				mcplib.Required(),
				mcplib.Description("Webex meeting ID — from webex_list_meetings or a meeting URL"),
			),
		),
		handleGetMeetingSummary,
	)

	s.AddTool(
		mcplib.NewTool("webex_list_teams",
			mcplib.WithDescription(
				"List Webex teams the authenticated user belongs to. "+
					"Teams are persistent collections of spaces and members — useful for "+
					"understanding org structure and locating team-specific spaces.",
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Maximum number of teams to return (1–200, default 20)"),
			),
		),
		handleListTeams,
	)
}

// ── Tool Handlers ─────────────────────────────────────────────────────────────

func handleListSpaces(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	keyword := strings.TrimSpace(req.GetString("keyword", ""))
	spaceType := strings.TrimSpace(req.GetString("space_type", ""))
	limit := req.GetInt("limit", 20)

	r := client.NewRequest(config.CallingBaseURL, "GET", "/rooms")
	r.QueryParam("max", strconv.Itoa(limit))
	if spaceType == "direct" || spaceType == "group" {
		r.QueryParam("type", spaceType)
	}
	body, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	spaces := parseItems(body)

	if keyword != "" {
		kw := strings.ToLower(keyword)
		filtered := spaces[:0]
		for _, s := range spaces {
			if strings.Contains(strings.ToLower(str(s["title"])), kw) {
				filtered = append(filtered, s)
			}
		}
		spaces = filtered
	}

	if len(spaces) == 0 {
		return jsonText(map[string]any{
			"count": 0, "spaces": []any{}, "message": "No matching spaces found.",
		}), nil
	}

	result := make([]map[string]any, 0, len(spaces))
	for _, s := range spaces {
		result = append(result, map[string]any{
			"id":           str(s["id"]),
			"title":        str(s["title"]),
			"type":         str(s["type"]),
			"lastActivity": fmtTS(str(s["lastActivity"])),
			"isLocked":     s["isLocked"],
		})
	}
	return jsonText(map[string]any{"count": len(result), "spaces": result}), nil
}

func handleGetMessages(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	roomID := strings.TrimSpace(req.GetString("room_id", ""))
	roomName := strings.TrimSpace(req.GetString("room_name", ""))
	limit := req.GetInt("limit", 20)

	if roomID == "" && roomName == "" {
		return mcplib.NewToolResultText("Error: Provide either room_id or room_name."), nil
	}

	// Name lookup — list rooms and filter client-side.
	if roomID == "" {
		r := client.NewRequest(config.CallingBaseURL, "GET", "/rooms")
		r.QueryParam("max", "200")
		listBody, _, err := r.Do()
		if err != nil {
			return mcplib.NewToolResultText("Error: " + err.Error()), nil
		}
		kw := strings.ToLower(roomName)
		for _, s := range parseItems(listBody) {
			if strings.Contains(strings.ToLower(str(s["title"])), kw) {
				roomID = str(s["id"])
				break
			}
		}
		if roomID == "" {
			return mcplib.NewToolResultText(
				fmt.Sprintf("Error: No space found matching '%s'.", roomName),
			), nil
		}
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/messages")
	r.QueryParam("roomId", roomID)
	r.QueryParam("max", strconv.Itoa(limit))
	msgBody, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	msgs := parseItems(msgBody)
	if len(msgs) == 0 {
		return jsonText(map[string]any{"count": 0, "room_id": roomID, "messages": []any{}}), nil
	}

	result := make([]map[string]any, 0, len(msgs))
	for _, m := range msgs {
		text := str(m["text"])
		if len(text) > 500 {
			text = text[:500]
		}
		result = append(result, map[string]any{
			"id":                str(m["id"]),
			"personEmail":       str(m["personEmail"]),
			"personDisplayName": str(m["personDisplayName"]),
			"text":              text,
			"created":           fmtTS(str(m["created"])),
			"hasFiles":          m["files"] != nil,
		})
	}
	return jsonText(map[string]any{
		"count": len(result), "room_id": roomID, "messages": result,
	}), nil
}

type msgMatch struct {
	SpaceTitle        string `json:"spaceTitle"`
	RoomID            string `json:"room_id"`
	PersonEmail       string `json:"personEmail"`
	PersonDisplayName string `json:"personDisplayName"`
	Text              string `json:"text"`
	Created           string `json:"created"`
}

func handleSearchMessages(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	keyword := strings.TrimSpace(req.GetString("keyword", ""))
	spaceLimit := req.GetInt("space_limit", msgSearchSpaceLimit)
	msgsPerSpace := req.GetInt("messages_per_space", 50)

	if keyword == "" {
		return mcplib.NewToolResultText("Error: keyword is required (min 2 chars)."), nil
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/rooms")
	r.QueryParam("max", strconv.Itoa(spaceLimit))
	listBody, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	spaces := parseItems(listBody)
	if len(spaces) == 0 {
		return jsonText(map[string]any{
			"keyword": keyword, "spaces_scanned": 0, "match_count": 0, "matches": []any{},
		}), nil
	}

	kwLower := strings.ToLower(keyword)
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		matches []msgMatch
	)

	for _, space := range spaces {
		space := space
		wg.Add(1)
		go func() {
			defer wg.Done()
			spaceID := str(space["id"])
			spaceTitle := str(space["title"])

			mr := client.NewRequest(config.CallingBaseURL, "GET", "/messages")
			mr.QueryParam("roomId", spaceID)
			mr.QueryParam("max", strconv.Itoa(msgsPerSpace))
			msgBody, _, err := mr.Do()
			if err != nil {
				return // best-effort; skip spaces that error
			}
			for _, m := range parseItems(msgBody) {
				text := str(m["text"])
				if strings.Contains(strings.ToLower(text), kwLower) {
					if len(text) > 500 {
						text = text[:500]
					}
					hit := msgMatch{
						SpaceTitle:        spaceTitle,
						RoomID:            spaceID,
						PersonEmail:       str(m["personEmail"]),
						PersonDisplayName: str(m["personDisplayName"]),
						Text:              text,
						Created:           fmtTS(str(m["created"])),
					}
					mu.Lock()
					matches = append(matches, hit)
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	if matches == nil {
		matches = []msgMatch{}
	}
	return jsonText(map[string]any{
		"keyword":        keyword,
		"spaces_scanned": len(spaces),
		"match_count":    len(matches),
		"matches":        matches,
	}), nil
}

func handleSendMessage(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	text := strings.TrimSpace(req.GetString("text", ""))
	roomID := strings.TrimSpace(req.GetString("room_id", ""))
	toPersonEmail := strings.TrimSpace(req.GetString("to_person_email", ""))
	markdown := strings.TrimSpace(req.GetString("markdown", ""))

	if text == "" {
		return mcplib.NewToolResultText("Error: text is required."), nil
	}
	if roomID == "" && toPersonEmail == "" {
		return mcplib.NewToolResultText("Error: Provide either room_id or to_person_email."), nil
	}

	r := client.NewRequest(config.CallingBaseURL, "POST", "/messages")
	r.BodyString("text", text)
	r.BodyString("roomId", roomID)
	r.BodyString("toPersonEmail", toPersonEmail)
	r.BodyString("markdown", markdown)
	body, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	var msg map[string]any
	if jsonErr := json.Unmarshal(body, &msg); jsonErr != nil || msg == nil {
		msg = map[string]any{}
	}
	return jsonText(map[string]any{
		"success":       true,
		"message_id":    str(msg["id"]),
		"created":       fmtTS(str(msg["created"])),
		"roomId":        str(msg["roomId"]),
		"toPersonEmail": str(msg["toPersonEmail"]),
	}), nil
}

func handleListDirectMessages(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	limit := req.GetInt("limit", 10)

	r := client.NewRequest(config.CallingBaseURL, "GET", "/rooms")
	r.QueryParam("type", "direct")
	r.QueryParam("max", strconv.Itoa(limit))
	listBody, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	spaces := parseItems(listBody)
	if len(spaces) == 0 {
		return jsonText(map[string]any{"count": 0, "conversations": []any{}}), nil
	}

	type lastMsg struct {
		From    string `json:"from"`
		Text    string `json:"text"`
		Created string `json:"created"`
	}
	type conversation struct {
		RoomID       string   `json:"room_id"`
		Title        string   `json:"title"`
		LastActivity string   `json:"lastActivity"`
		LastMessage  *lastMsg `json:"last_message"`
	}

	convs := make([]conversation, len(spaces))
	var wg sync.WaitGroup

	for i, space := range spaces {
		i, space := i, space
		convs[i] = conversation{
			RoomID:       str(space["id"]),
			Title:        str(space["title"]),
			LastActivity: fmtTS(str(space["lastActivity"])),
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			mr := client.NewRequest(config.CallingBaseURL, "GET", "/messages")
			mr.QueryParam("roomId", str(space["id"]))
			mr.QueryParam("max", "1")
			msgBody, _, err := mr.Do()
			if err != nil {
				return
			}
			msgs := parseItems(msgBody)
			if len(msgs) == 0 {
				return
			}
			m := msgs[0]
			text := str(m["text"])
			if len(text) > 200 {
				text = text[:200]
			}
			convs[i].LastMessage = &lastMsg{
				From:    str(m["personEmail"]),
				Text:    text,
				Created: fmtTS(str(m["created"])),
			}
		}()
	}
	wg.Wait()

	return jsonText(map[string]any{"count": len(convs), "conversations": convs}), nil
}

func handleGetPerson(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	email := strings.TrimSpace(req.GetString("email", ""))
	displayName := strings.TrimSpace(req.GetString("display_name", ""))

	if email == "" && displayName == "" {
		return mcplib.NewToolResultText("Error: Provide either email or display_name."), nil
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/people")
	if email != "" {
		r.QueryParam("email", email)
	} else {
		r.QueryParam("displayName", displayName)
	}
	body, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	people := parseItems(body)
	if len(people) == 0 {
		query := email
		if query == "" {
			query = displayName
		}
		return jsonText(map[string]any{
			"count": 0, "message": fmt.Sprintf("No person found matching '%s'.", query),
		}), nil
	}

	result := make([]map[string]any, 0, len(people))
	for _, p := range people {
		result = append(result, map[string]any{
			"id":           str(p["id"]),
			"displayName":  str(p["displayName"]),
			"emails":       p["emails"],
			"firstName":    str(p["firstName"]),
			"lastName":     str(p["lastName"]),
			"title":        str(p["title"]),
			"orgId":        str(p["orgId"]),
			"status":       str(p["status"]),
			"avatar":       str(p["avatar"]),
			"lastActivity": fmtTS(str(p["lastActivity"])),
		})
	}
	return jsonText(map[string]any{"count": len(result), "people": result}), nil
}

func handleListMeetings(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	limit := req.GetInt("limit", 20)
	now := time.Now().UTC()

	fromDate := strings.TrimSpace(req.GetString("from_date", ""))
	if fromDate == "" {
		fromDate = now.AddDate(0, 0, -7).Truncate(24 * time.Hour).Format(time.RFC3339)
	}
	toDate := strings.TrimSpace(req.GetString("to_date", ""))
	if toDate == "" {
		end := now.AddDate(0, 0, 30)
		toDate = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/meetings")
	r.QueryParam("from", fromDate)
	r.QueryParam("to", toDate)
	r.QueryParam("max", strconv.Itoa(limit))
	body, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	meetings := parseItems(body)
	if len(meetings) == 0 {
		return jsonText(map[string]any{
			"count": 0, "from": fromDate, "to": toDate, "meetings": []any{},
		}), nil
	}

	result := make([]map[string]any, 0, len(meetings))
	for _, m := range meetings {
		result = append(result, map[string]any{
			"id":            str(m["id"]),
			"title":         str(m["title"]),
			"start":         fmtTS(str(m["start"])),
			"end":           fmtTS(str(m["end"])),
			"hostEmail":     str(m["hostEmail"]),
			"status":        str(m["status"]),
			"meetingNumber": str(m["meetingNumber"]),
			"webLink":       str(m["webLink"]),
		})
	}
	return jsonText(map[string]any{
		"count": len(result), "from": fromDate, "to": toDate, "meetings": result,
	}), nil
}

func handleGetMeetingTranscript(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	meetingID, err := req.RequireString("meeting_id")
	if err != nil {
		return mcplib.NewToolResultText("Error: meeting_id is required."), nil
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/meetingTranscripts")
	r.QueryParam("meetingId", meetingID)
	body, _, apiErr := r.Do()
	if apiErr != nil {
		return mcplib.NewToolResultText("Error: " + apiErr.Error()), nil
	}

	transcripts := parseItems(body)
	if len(transcripts) == 0 {
		return jsonText(map[string]any{
			"meeting_id":  meetingID,
			"message":     "No transcript found. Transcription may not have been enabled for this meeting, or it may still be processing.",
			"transcripts": []any{},
		}), nil
	}

	result := make([]map[string]any, 0, len(transcripts))
	for _, t := range transcripts {
		result = append(result, map[string]any{
			"id":        str(t["id"]),
			"meetingId": str(t["meetingId"]),
			"status":    str(t["status"]),
			"created":   fmtTS(str(t["created"])),
		})
	}
	return jsonText(map[string]any{
		"meeting_id":  meetingID,
		"item_count":  len(result),
		"transcripts": result,
	}), nil
}

func handleGetMeetingSummary(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	meetingID, err := req.RequireString("meeting_id")
	if err != nil {
		return mcplib.NewToolResultText("Error: meeting_id is required."), nil
	}

	r := client.NewRequest(config.CallingBaseURL, "GET", "/meetingSummaries")
	r.QueryParam("meetingId", meetingID)
	body, _, apiErr := r.Do()
	if apiErr != nil {
		return mcplib.NewToolResultText("Error: " + apiErr.Error()), nil
	}

	summaries := parseItems(body)
	if len(summaries) == 0 {
		return jsonText(map[string]any{
			"meeting_id": meetingID,
			"message":    "No AI summary found. The meeting may still be processing, or AI summaries may be disabled for your org.",
			"summaries":  []any{},
		}), nil
	}

	coalesceList := func(v any) any {
		if v == nil {
			return []any{}
		}
		return v
	}

	result := make([]map[string]any, 0, len(summaries))
	for _, s := range summaries {
		result = append(result, map[string]any{
			"id":          str(s["id"]),
			"meetingId":   str(s["meetingId"]),
			"keyPoints":   coalesceList(s["keyPoints"]),
			"actionItems": coalesceList(s["actionItems"]),
			"decisions":   coalesceList(s["decisions"]),
			"created":     fmtTS(str(s["created"])),
		})
	}
	return jsonText(map[string]any{
		"meeting_id": meetingID,
		"item_count": len(result),
		"summaries":  result,
	}), nil
}

func handleListTeams(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	limit := req.GetInt("limit", 20)

	r := client.NewRequest(config.CallingBaseURL, "GET", "/teams")
	r.QueryParam("max", strconv.Itoa(limit))
	body, _, err := r.Do()
	if err != nil {
		return mcplib.NewToolResultText("Error: " + err.Error()), nil
	}

	teams := parseItems(body)
	if len(teams) == 0 {
		return jsonText(map[string]any{"count": 0, "teams": []any{}}), nil
	}

	result := make([]map[string]any, 0, len(teams))
	for _, t := range teams {
		result = append(result, map[string]any{
			"id":        str(t["id"]),
			"name":      str(t["name"]),
			"creatorId": str(t["creatorId"]),
			"created":   fmtTS(str(t["created"])),
		})
	}
	return jsonText(map[string]any{"count": len(result), "teams": result}), nil
}
