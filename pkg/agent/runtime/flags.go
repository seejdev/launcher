package runtime

import (
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
