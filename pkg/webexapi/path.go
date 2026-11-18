package webexapi

import (
	"net/url"
	"strings"
)

// FillPath replaces each {name} in a path template with its value, escaped for a
// path segment. A parameter with an empty value is left as it is.
func FillPath(template string, params map[string]string) string {
	for k, v := range params {
		if v != "" {
			template = strings.ReplaceAll(template, "{"+k+"}", url.PathEscape(v))
		}
	}
	return template
}
