package main

import "os"

type AppConfig struct {
	DatabaseDSN string `env:"DATABASE_URL"`
	ServerAddr  string `env:"SERVER_ADDR"`
}

func NewAppConfig() (*AppConfig, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, ErrMissingEnvVar("DATABASE_URL")
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8081"
	}

	return &AppConfig{
		DatabaseDSN: dbURL,
		ServerAddr:  serverAddr,
	}, nil
}

func ErrMissingEnvVar(name string) error {
	return &MissingEnvVarError{Name: name}
}

type MissingEnvVarError struct {
	Name string
}

func (e *MissingEnvVarError) Error() string {
	return "missing env var " + e.Name
}
