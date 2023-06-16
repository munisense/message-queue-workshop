package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	MQProtocol string
	MQHost     string
	MQVHost    string
	MQPort     int
	MQUsername string
	MQPassword string
}

func LoadConfig() *Config {
	// Load environment files from .env if available
	_ = godotenv.Load(".env")

	if err := godotenv.Load(".env.default"); err != nil {
		panic(err)
	}

	MQPort, err := strconv.Atoi(os.Getenv("MQ_PORT"))
	if err != nil {
		panic("failed to parse MQ_PORT as integer")
	}

	c := &Config{
		MQProtocol: os.Getenv("MQ_PROTOCOL"),
		MQHost:     os.Getenv("MQ_HOST"),
		MQVHost:    os.Getenv("MQ_VHOST"),
		MQPort:     MQPort,
		MQUsername: os.Getenv("MQ_USER"),
		MQPassword: os.Getenv("MQ_PWD"),
	}

	if c.MQUsername == "" {
		panic("Please configure the MQ_USER environment variable in your .env file")
	}

	if c.MQPassword == "" {
		panic("Please configure the MQ_PWD environment variable in your .env file")
	}

	if c.MQVHost == "" {
		panic("Please configure the MQ_VHOST environment variable in your .env file")
	}

	if c.MQHost == "" {
		panic("Please configure the MQ_HOST environment variable in your .env file")
	}

	return c
}
