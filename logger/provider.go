package logger

import (
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	ProvideLogger,
)
