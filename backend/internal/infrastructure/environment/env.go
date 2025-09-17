package environment

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"time"
)

func ProvideEnvironment() *Environment {
	e := &Environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	return e
}

type Environment struct {
	Env string `env:"ENV,required"`
	ServerEnvironment
	VectorDatabaseEnvironment
	AppDatabaseEnvironment
	GoogleCloudEnvironment
	CloudStorageEnvironment
	DocumentAIEnvironment
	VertexAIEnvironment
	SyncQueueEnvironment
	RedisEnvironment
}

type ServerEnvironment struct {
	ListenAddress   string        `env:"LISTEN_ADDRESS,required"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
}

type VectorDatabaseEnvironment struct {
	Host     string `env:"VECTOR_DB_HOST,required"`
	Port     int    `env:"VECTOR_DB_PORT,required"`
	Username string `env:"VECTOR_DB_USERNAME,required"`
	Password string `env:"VECTOR_DB_PASSWORD,required"`
	Database string `env:"VECTOR_DB_NAME,required"`
	SSLMode  string `env:"VECTOR_DB_SSL_MODE,required"`
}

func (v *VectorDatabaseEnvironment) VectorDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		v.Username, v.Password, v.Host, v.Port, v.Database, v.SSLMode)
}

type AppDatabaseEnvironment struct {
	Host     string `env:"APP_DB_HOST,required"`
	Port     int    `env:"APP_DB_PORT,required"`
	Username string `env:"APP_DB_USERNAME,required"`
	Password string `env:"APP_DB_PASSWORD,required"`
	Database string `env:"APP_DB_NAME,required"`
	SSLMode  string `env:"APP_DB_SSL_MODE,required"`
}

func (a *AppDatabaseEnvironment) AppDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		a.Username, a.Password, a.Host, a.Port, a.Database, a.SSLMode)
}

type GoogleCloudEnvironment struct {
	ProjectID string `env:"GOOGLE_CLOUD_PROJECT_ID,required"`
}

type CloudStorageEnvironment struct {
	BucketName string `env:"CLOUD_STORAGE_BUCKET_NAME,required"`
}

type DocumentAIEnvironment struct {
	DocumentAILocation string `env:"DOCUMENT_AI_LOCATION,required"`
	ProcessorID        string `env:"DOCUMENT_AI_PROCESSOR_ID,required"`
}

type VertexAIEnvironment struct {
	VertexAILocation string `env:"VERTEX_AI_LOCATION,required"`
}

type SyncQueueEnvironment struct {
	QueueName     string `env:"SYNC_QUEUE_NAME"`
	QueueLocation string `env:"SYNC_QUEUE_LOCATION"`
	TargetURL     string `env:"SYNC_TARGET_URL"`
}

type RedisEnvironment struct {
	RedisURL string `env:"REDIS_URL,required"`
}
