package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// SDK returns a transport for the SDK packages, set up from this run's flags like
// every other request: the token and its refresh, 429 retries, debug output,
// read-only mode, dry-run and Contact Center pacing.
func SDK() *webexapi.Client {
	api := &webexapi.Client{
		HTTP: http.DefaultClient,
		Token: func(context.Context) (string, error) {
			return config.Token(), nil
		},
		MaxRetries:   config.MaxRetry(),
		MaxRetryWait: time.Duration(config.MaxRetryTimer()) * time.Second,
		OnRetry: func(wait time.Duration, attempt, max int) {
			fmt.Fprintf(os.Stderr, "Rate limited (429). Retrying in %s (attempt %d/%d)...\n", wait.Round(time.Second), attempt, max)
		},
		BeforeSend: sdkBeforeSend,
		AfterResponse: func(status, size int) {
			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", status, size)
			}
		},
	}
	if config.TokenRefresher != nil {
		api.Refresh = func(context.Context) error {
			newToken, err := config.TokenRefresher()
			if err != nil {
				return err
			}
			config.SetToken(newToken)
			return nil
		}
	}
	return api
}

// sdkBeforeSend applies debug output, read-only mode and dry-run to each attempt.
// SDK requests carry no route template, so read-only mode allows reads only.
func sdkBeforeSend(r *http.Request) error {
	body := ""
	if r.GetBody != nil {
		if rc, err := r.GetBody(); err == nil {
			data, _ := io.ReadAll(rc)
			body = string(data)
		}
	}
	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", r.Method, r.URL)
		for k, v := range r.Header {
			fmt.Fprintf(os.Stderr, "DEBUG:   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
		}
		if body != "" {
			fmt.Fprintf(os.Stderr, "DEBUG:   Body: %s\n", truncate(body, 500))
		}
	}
	if config.ReadOnly() && isWriteMethod(r.Method) {
		return ReadOnlyError(r.Method, r.URL.String())
	}
	if config.DryRun() && isWriteMethod(r.Method) {
		fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", r.Method, r.URL)
		if body != "" {
			fmt.Fprintf(os.Stderr, "[DRY RUN]   Body: %s\n", body)
		}
		return ErrDryRun
	}
	return ccPacer().Wait(r)
}
