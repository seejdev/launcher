package runtime

import (
	"time"

	"github.com/go-kit/kit/log"
)

type Flags struct {
	logger log.Logger
	// flags
}

func NewFlags(logger log.Logger) *Flags {
	f := &Flags{
		logger: logger,
	}

	return f
}

func SetDesktopEnabled(enabled bool) error {
	return nil // TODO
}

func DesktopEnabled() bool {
	return true // TODO
}

func SetControlRequestInterval() error {
	return nil // TODO
}

func ControlRequestInterval() time.Duration {
	return 60 * time.Second // TODO
}


// DisableControlTLS disables TLS transport with the control server.
DisableControlTLS bool
// InsecureControlTLS disables TLS certificate validation for the control server.
InsecureControlTLS bool