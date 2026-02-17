package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// Do executes an HTTP request. On a 401, it attempts to refresh the token and retry once.
func Do(req *Request) ([]byte, int, error) {
	body, status, err := doOnce(req)
	if status == 401 && config.TokenRefresher != nil {
		newToken, refreshErr := config.TokenRefresher()
		if refreshErr != nil {
			return body, status, err
		}
		config.SetToken(newToken)
		return doOnce(req)
	}
	return body, status, err
}

// doOnce executes a single HTTP request without retry.
func doOnce(req *Request) ([]byte, int, error) {
	// Build URL
	url := req.baseURL + req.path
	for k, v := range req.pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", v)
	}

	// Add query params
	if len(req.queryParams) > 0 {
		parts := make([]string, 0, len(req.queryParams))
		for k, v := range req.queryParams {
			if v != "" {
				parts = append(parts, fmt.Sprintf("%s=%s", k, v))
			}
		}
		if len(parts) > 0 {
			url += "?" + strings.Join(parts, "&")
		}
	}

	// Build body
	var bodyReader io.Reader
	if req.bodyRaw != "" {
		bodyReader = strings.NewReader(req.bodyRaw)
	}

	httpReq, err := http.NewRequest(req.method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("building request: %w", err)
	}

	// Auth
	httpReq.Header.Set("Authorization", "Bearer "+config.Token())

	// Content type
	if req.bodyRaw != "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	// Extra headers
	for k, v := range req.headers {
		httpReq.Header.Set(k, v)
	}

	// Debug
	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", req.method, url)
		for k, v := range httpReq.Header {
			if k != "Authorization" {
				fmt.Fprintf(os.Stderr, "DEBUG:   %s: %s\n", k, strings.Join(v, ", "))
			}
		}
		if req.bodyRaw != "" {
			fmt.Fprintf(os.Stderr, "DEBUG:   Body: %s\n", truncate(req.bodyRaw, 500))
		}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", resp.StatusCode, len(body))
	}

	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, fmt.Errorf("API error %d: %s", resp.StatusCode, truncate(string(body), 500))
	}

	return body, resp.StatusCode, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
