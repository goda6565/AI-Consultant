package entity

import (
	sharedValue "github.com/goda6565/ai-consultant/backend/internal/domain/shared/value"
)

type JobConfig struct {
	id                   sharedValue.ID
	problemID            sharedValue.ID
	enableInternalSearch bool
}

func NewJobConfig(id sharedValue.ID, problemID sharedValue.ID, enableInternalSearch bool) *JobConfig {
	return &JobConfig{id: id, problemID: problemID, enableInternalSearch: enableInternalSearch}
}

func (j *JobConfig) GetID() sharedValue.ID {
	return j.id
}

func (j *JobConfig) GetProblemID() sharedValue.ID {
	return j.problemID
}

func (j *JobConfig) GetEnableInternalSearch() bool {
	return j.enableInternalSearch
}

func (j *JobConfig) EnableInternalSearch() {
	j.enableInternalSearch = true
}

func (j *JobConfig) DisableInternalSearch() {
	j.enableInternalSearch = false
}
