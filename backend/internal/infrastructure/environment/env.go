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
	Env string `env:"ENV" envDefault:"development"`
	ServerEnvironment
	VectorDatabaseEnvironment
	AppDatabaseEnvironment
	GoogleCloudEnvironment
	CloudStorageEnvironment
	DocumentAIEnvironment
	VertexAIEnvironment
}

type ServerEnvironment struct {
	ListenAddress   string        `env:"LISTEN_ADDRESS" envDefault:"localhost:8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type VectorDatabaseEnvironment struct {
	Host     string `env:"VECTOR_DB_HOST" envDefault:"localhost"`
	Port     int    `env:"VECTOR_DB_PORT" envDefault:"5432"`
	Username string `env:"VECTOR_DB_USERNAME" envDefault:"postgres"`
	Password string `env:"VECTOR_DB_PASSWORD" envDefault:"password"`
	Database string `env:"VECTOR_DB_NAME" envDefault:"vector_db"`
	SSLMode  string `env:"VECTOR_DB_SSL_MODE" envDefault:"disable"`
}

func (v *VectorDatabaseEnvironment) VectorDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		v.Username, v.Password, v.Host, v.Port, v.Database, v.SSLMode)
}

type AppDatabaseEnvironment struct {
	Host     string `env:"APP_DB_HOST" envDefault:"localhost"`
	Port     int    `env:"APP_DB_PORT" envDefault:"5432"`
	Username string `env:"APP_DB_USERNAME" envDefault:"postgres"`
	Password string `env:"APP_DB_PASSWORD" envDefault:"password"`
	Database string `env:"APP_DB_NAME" envDefault:"app_db"`
	SSLMode  string `env:"APP_DB_SSL_MODE" envDefault:"disable"`
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
	DocumentAILocation string `env:"DOCUMENT_AI_LOCATION" envDefault:"us"`
	ProcessorID        string `env:"DOCUMENT_AI_PROCESSOR_ID,required"`
}

type VertexAIEnvironment struct {
	VertexAILocation string `env:"VERTEX_AI_LOCATION" envDefault:"global"`
}
