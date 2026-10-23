package auth

// HeadlessReason returns why this environment probably has no browser for the
// OAuth redirect to reach, or "" when the browser flow should work.
func HeadlessReason(getenv func(string) string, goos string) string {
	switch {
	case getenv("SSH_CONNECTION") != "" || getenv("SSH_TTY") != "":
		return "SSH session"
	case getenv("CI") != "":
		return "CI environment"
	case goos == "linux" && getenv("DISPLAY") == "" && getenv("WAYLAND_DISPLAY") == "":
		return "no graphical display"
	}
	return ""
}
