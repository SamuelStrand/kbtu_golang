package modules

import "time"

type PostgreConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	DBName      string
	SSLMode     string
	ExecTimeout time.Duration
}

type AppConfig struct {
	ServerAddr string
	APIKey     string
	Postgres   PostgreConfig
}
