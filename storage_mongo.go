package storage

import (
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	// Default connection settings
	Hostname     string
	Hostnames    []string
	Port         uint16 // 0 to 65,535
	DatabaseName string

	// Credential Access
	Username string
	Password string

	// Replica Set
	Replset string
}

// connectionURL generates a formatted Mongo Connection URL
func ConnectionURI(cfg *Config) string {
	connectionString := "mongodb://"
	credentials := ""

	if cfg.Username != "" && cfg.Password != "" {
		credentials = fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
	}

	hosts := ""
	if cfg.Hostnames != nil && cfg.Hostname == "" {
		hosts = strings.Join(cfg.Hostnames, fmt.Sprintf(":%s,", strconv.Itoa(int(cfg.Port))))
		cfg.Hostname = hosts
	}

	connectionString = fmt.Sprintf("%s%s%s:%s/%s",
		connectionString,
		credentials,
		cfg.Hostname,
		strconv.Itoa(int(cfg.Port)),
		cfg.DatabaseName,
	)

	if cfg.Replset != "" {
		connectionString += "?replicaSet=" + cfg.Replset
	}

	return connectionString
}

// NewDatastore returns a connection to your configured MongoDB
func NewDatastore(cfg *Config) string {
	return ConnectionURI(cfg)
}
