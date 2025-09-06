package value

import "fmt"

type StoragePath struct {
	rawPath    string
	bucketName string
	objectName string
}

func (p StoragePath) Equals(other StoragePath) bool {
	if p.rawPath != other.rawPath {
		return false
	}
	if p.bucketName != other.bucketName {
		return false
	}
	if p.objectName != other.objectName {
		return false
	}
	return true
}

func (p StoragePath) RawPath() string {
	return p.rawPath
}

func (p StoragePath) BucketName() string {
	return p.bucketName
}

func (p StoragePath) ObjectName() string {
	return p.objectName
}

// TODO: validate path
func NewStoragePath(bucketName string, objectName string) StoragePath {
	return StoragePath{
		rawPath:    fmt.Sprintf("gs://%s/%s", bucketName, objectName),
		bucketName: bucketName,
		objectName: objectName,
	}
}
