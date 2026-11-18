package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/readonly"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// ErrDryRun is returned when a write operation is intercepted by --dry-run mode.
var ErrDryRun = errors.New("dry run: no changes made")

// ErrReadOnly is returned when read-only mode blocks a request.
var ErrReadOnly = errors.New("read-only mode")

// ReadOnlyError reports a request blocked by read-only mode.
func ReadOnlyError(method, url string) error {
	return fmt.Errorf("%w: %s %s not sent — run 'webex login' to leave read-only mode", ErrReadOnly, method, url)
}

// Do executes an HTTP request.
// On a 401 it refreshes the token and retries once.
// On a 429 it reads Retry-After and retries up to 3 times.
func Do(req *Request) ([]byte, int, error) {
	resp, err := API(req.path, req.logBody()).Send(context.Background(), req.sdkRequest())
	if resp == nil {
		return nil, 0, CLIError(err)
	}
	return resp.Body, resp.Status, CLIError(err)
}

// API is a client set up by the command line: its token and refresh, 429 retry
// limits, --debug output, read-only mode and --dry-run. template is the request's
// unexpanded path, which is how read-only mode lists query endpoints, and logBody
// stands in for the body in diagnostics.
func API(template, logBody string) *webexapi.Client {
	c := &webexapi.Client{
		Origins:      config.Origins(),
		Token:        func(context.Context) (string, error) { return config.Token(), nil },
		MaxRetries:   config.MaxRetry(),
		MaxRetryWait: time.Duration(config.MaxRetryTimer()) * time.Second,
		BeforeSend: func(httpReq *http.Request) error {
			url := httpReq.URL.String()
			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", httpReq.Method, url)
				for k, v := range httpReq.Header {
					fmt.Fprintf(os.Stderr, "DEBUG:   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
				}
				if logBody != "" {
					fmt.Fprintf(os.Stderr, "DEBUG:   Body: %s\n", truncate(logBody, 500))
				}
			}
			if config.ReadOnly() && !readOnlyAllows(httpReq.Method, template) {
				return ReadOnlyError(httpReq.Method, url)
			}
			// Dry-run: intercept write operations before making the HTTP call
			if config.DryRun() && isWriteMethod(httpReq.Method) {
				fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", httpReq.Method, url)
				for k, v := range httpReq.Header {
					fmt.Fprintf(os.Stderr, "[DRY RUN]   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
				}
				if logBody != "" {
					fmt.Fprintf(os.Stderr, "[DRY RUN]   Body: %s\n", logBody)
				}
				return ErrDryRun
			}
			return nil
		},
		AfterResponse: func(status, size int) {
			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", status, size)
			}
		},
		OnRetry: func(wait time.Duration, attempt, max int) {
			fmt.Fprintf(os.Stderr, "Rate limited (429). Retrying in %s (attempt %d/%d)...\n", wait.Round(time.Second), attempt, max)
		},
	}
	if refresh := config.TokenRefresher; refresh != nil {
		c.Refresh = func(context.Context) error {
			token, err := refresh()
			if err != nil {
				return err
			}
			config.SetToken(token)
			return nil
		}
	}
	return c
}

// CLIError says how to fix a request sent to the wrong data region, keeping the
// context of any error that wraps it, and returns any other error as it is.
func CLIError(err error) error {
	var wrongRegion *webexapi.WrongRegionError
	if !errors.As(err, &wrongRegion) {
		return err
	}
	msg := fmt.Sprintf("wrong data region (451): this organization's data is not hosted in region %q", regionLabel())
	if wrongRegion.Endpoint != "" {
		msg += fmt.Sprintf("; Webex says to use %s", wrongRegion.Endpoint)
	}
	msg = fmt.Sprintf("%s — rerun with --region <%s> (or: webex config set region <region>)",
		msg, strings.Join(config.Regions, "|"))
	return &cliError{msg: strings.Replace(err.Error(), wrongRegion.Error(), msg, 1), err: err}
}

// cliError is an error with the CLI's own message.
type cliError struct {
	msg string
	err error
}

func (e *cliError) Error() string { return e.msg }
func (e *cliError) Unwrap() error { return e.err }

func regionLabel() string {
	if r := config.Region(); r != "" {
		return r
	}
	return "us (default)"
}

// redactHeader replaces a credential value with its length, which is enough to
// tell a missing or truncated token from a present one without printing it.
func redactHeader(key, value string) string {
	// An AI agent's runtime token lets anyone talk to that agent.
	if strings.EqualFold(key, "Bot-Access-Token") {
		return fmt.Sprintf("<%d chars>", len(value))
	}
	if !strings.EqualFold(key, "Authorization") {
		return value
	}
	if scheme, cred, ok := strings.Cut(value, " "); ok {
		return fmt.Sprintf("%s <%d chars>", scheme, len(cred))
	}
	return fmt.Sprintf("<%d chars>", len(value))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// readOnlyAllows reports whether a request may be sent in read-only mode.
// path is the unexpanded template, which is how query endpoints are listed.
func readOnlyAllows(method, path string) bool {
	if method == http.MethodPost {
		return readonly.IsQueryPOST(path)
	}
	return !isWriteMethod(method)
}

// isWriteMethod returns true for HTTP methods that modify data.
func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}
