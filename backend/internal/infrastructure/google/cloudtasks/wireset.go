package cloudtasks

import "github.com/google/wire"

var Set = wire.NewSet(
	NewCloudTasksClient,
)
