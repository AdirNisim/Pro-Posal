package config

import (
	"os"
	"strconv"
)

var AppConfig Config

const DEFAULT_AUTH_EXPIRATION_TIME_MINUTES = "5"

type Config struct {
	Database Database
	Server   Server
	Auth     Auth
}

type Server struct {
	Port string
}

type Auth struct {
	ExpirationTimeMinutes int
	JWTSigningSecret      string
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
	AppConfig.Auth.loadConfig()
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

func (a *Auth) loadConfig() {
	expirationTimeStr := getValueOrDefault("AUTH_EXPIRATION_TIME_MIN", DEFAULT_AUTH_EXPIRATION_TIME_MINUTES)
	expirationTime, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		panic("Invalid AUTH_EXPIRATION_TIME_MIN")
	}
	a.ExpirationTimeMinutes = expirationTime

	a.JWTSigningSecret = os.Getenv("JWT_SIGNING_SECRET")
}

func getValueOrDefault(keyName string, defaultValue string) string {
	value := os.Getenv(keyName)
	if value == "" {
		return defaultValue
	}
	return value
}
