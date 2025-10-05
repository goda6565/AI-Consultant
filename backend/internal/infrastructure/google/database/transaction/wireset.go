package transaction

import "github.com/google/wire"

var Set = wire.NewSet(
	NewVectorUnitOfWork,
	NewAdminUnitOfWork,
)
