package cloudrunjob

import "github.com/google/wire"

var Set = wire.NewSet(
	NewCloudRunJobClient,
)
