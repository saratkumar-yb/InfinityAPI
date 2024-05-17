package config

import (
	"log"

	"github.com/go-ini/ini"
)

type Config struct {
	DBHost       string
	DBPort       int
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	HTTPListener string
	HTTPPort     int
}

var AppConfig Config

func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	dbSection := cfg.Section("db")
	serverSection := cfg.Section("server")

	AppConfig = Config{
		DBHost:       dbSection.Key("host").String(),
		DBPort:       dbSection.Key("port").MustInt(),
		DBUser:       dbSection.Key("user").String(),
		DBPassword:   dbSection.Key("password").String(),
		DBName:       dbSection.Key("dbname").String(),
		DBSSLMode:    dbSection.Key("sslmode").String(),
		HTTPListener: serverSection.Key("http_listener").String(),
		HTTPPort:     serverSection.Key("http_port").MustInt(),
	}
}
