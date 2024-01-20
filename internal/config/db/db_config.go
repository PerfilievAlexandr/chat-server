package config

import (
	configInterface "chat-server/internal/config/interface"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

var _ configInterface.DatabaseConfig = (*dbConfig)(nil)

const (
	port         = "DB_PORT"
	host         = "DB_HOST"
	driver       = "DB_DRIVER"
	user         = "DB_USER"
	password     = "DB_PASSWORD"
	sslMode      = "DB_SSL_MODE"
	databaseName = "DB_NAME"
)

type dbConfig struct {
	port         string
	host         string
	driver       string
	user         string
	password     string
	sslMode      string
	databaseName string
}

func NewDbConfig() (configInterface.DatabaseConfig, error) {
	port := os.Getenv(port)
	if len(port) == 0 {
		return nil, errors.New("db port not found")
	}

	host := os.Getenv(host)
	if len(host) == 0 {
		return nil, errors.New("db host not found")
	}

	driver := os.Getenv(driver)
	if len(driver) == 0 {
		return nil, errors.New("db driver not found")
	}

	user := os.Getenv(user)
	if len(user) == 0 {
		return nil, errors.New("db user not found")
	}

	password := os.Getenv(password)
	if len(password) == 0 {
		return nil, errors.New("db password not found")
	}

	sslMode := os.Getenv(sslMode)
	if len(sslMode) == 0 {
		return nil, errors.New("db sslMode not found")
	}

	databaseName := os.Getenv(databaseName)
	if len(databaseName) == 0 {
		return nil, errors.New("db databaseName not found")
	}

	return &dbConfig{
		host:         host,
		port:         port,
		driver:       driver,
		user:         user,
		password:     password,
		sslMode:      sslMode,
		databaseName: databaseName,
	}, nil
}

func (s *dbConfig) ConnectString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		s.host,
		s.port,
		s.user,
		s.password,
		s.databaseName,
		s.sslMode,
	)
}
