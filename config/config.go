package config

import (
	"os"
)

var AppConfig Config

type Config struct {
	Database Database
	Server   Server
}

type Server struct {
	Port string
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	SSLMode  string
}

func init() {
	AppConfig.Server.loadConfig()
	AppConfig.Database.loadConfig()
}

func (s *Server) loadConfig() {
	s.Port = os.Getenv("SERVER_PORT")
}

func (d *Database) loadConfig() {
	d.User = os.Getenv("DB_USER")
	d.Password = os.Getenv("DB_PASS")
	d.Host = os.Getenv("DB_HOST")
	d.Port = os.Getenv("DB_PORT")
	d.Name = os.Getenv("DB_NAME")
	d.SSLMode = os.Getenv("DB_SSLMODE")
}
