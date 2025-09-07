package value

import (
	"fmt"
	"unicode/utf8"

	"github.com/goda6565/ai-consultant/backend/internal/domain/errors"
)

const maxBucketNameLength = 50
const maxObjectNameLength = 50

type StorageInfo struct {
	bucketName string
	objectName string
}

func (p StorageInfo) Equals(other StorageInfo) bool {
	if p.bucketName != other.bucketName {
		return false
	}
	if p.objectName != other.objectName {
		return false
	}
	return true
}

func (p StorageInfo) BucketName() string {
	return p.bucketName
}

func (p StorageInfo) ObjectName() string {
	return p.objectName
}

func (p StorageInfo) Validate() error {
	if utf8.RuneCountInString(p.bucketName) > maxBucketNameLength {
		return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("bucket name must be less than %d characters", maxBucketNameLength))
	}
	if utf8.RuneCountInString(p.objectName) > maxObjectNameLength {
		return errors.NewDomainError(errors.ValidationError, fmt.Sprintf("object name must be less than %d characters", maxObjectNameLength))
	}
	return nil
}
func NewStorageInfo(bucketName string, objectName string) StorageInfo {
	return StorageInfo{
		bucketName: bucketName,
		objectName: objectName,
	}
}
