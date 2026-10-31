package auth

import (
	"os"
	"testing"
	"time"
)

// Every test gets a scratch credentials directory so none touches the real one.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "webex-auth-test-")
	if err != nil {
		panic(err)
	}
	os.Setenv("XDG_CONFIG_HOME", dir)
	os.Unsetenv(TokenStoreEnv)
	os.Unsetenv(RefreshTokenEnv)
	pollUnit = time.Millisecond
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}
