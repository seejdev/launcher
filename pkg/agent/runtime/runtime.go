package runtime

import (
	"github.com/go-kit/kit/log"
)

type Runtime struct {
	logger log.Logger
	Flags  *Flags
	// flags
	// Querier
	// storage?
	// other context stuff?
}

func NewRuntime(logger log.Logger, flags *Flags) *Runtime {
	r := &Runtime{
		logger: logger,
		Flags:  flags,
	}

	return r
}
