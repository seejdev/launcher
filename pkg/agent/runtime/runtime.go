package runtime

import (
	"github.com/go-kit/kit/log"
)

type Runtime struct {
	logger log.Logger
	Flags  Flags
	// flags
	// Querier
	// storage?
	// other context stuff?
}

func New(logger log.Logger) *Runtime {
	r := &Runtime{
		logger: logger,
	}

	return r
}
