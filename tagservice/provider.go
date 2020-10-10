package tagservice

import (
	"github.com/google/wire"
)

var SuperSet = wire.NewSet(
	ProvideDB,
	ProvideTagQueryRepository,
	ProvideTagServer,
)
