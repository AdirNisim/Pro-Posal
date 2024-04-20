package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

var AppConfig Config

func init() {
	err := godotenv.Load("../_.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	AppConfig = Config{
		Port: os.Getenv("PORT"),
	}
}
