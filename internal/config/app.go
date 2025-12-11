package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	S3       S3Config
	DynamoDB DynamoDBConfig
	SNS      SNSConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type S3Config struct {
	Endpoint     string
	AccessKey    string
	SecretKey    string
	BucketName   string
	Region       string
	UsePathStyle bool
}

type DynamoDBConfig struct {
	Endpoint  string
	TableName string
	Region    string
}

type SNSConfig struct {
	Mock     bool
	TopicARN string
	Region   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "aws_segundaentrega"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		S3: S3Config{
			Endpoint:     getEnv("S3_ENDPOINT", "http://localhost:9000"),
			AccessKey:    getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey:    getEnv("S3_SECRET_KEY", "minioadmin"),
			BucketName:   getEnv("S3_BUCKET_NAME", "aws-segundaentrega"),
			Region:       getEnv("S3_REGION", "us_east-1"),
			UsePathStyle: getEnv("S3_USE_PATH_STYLE", "true") == "true",
		},
		DynamoDB: DynamoDBConfig{
			Endpoint:  getEnv("DYNAMODB_ENDPOINT", "http://localhost:8000"),
			TableName: getEnv("DYNAMODB_TABLE", "sesiones-alumnos"),
			Region:    getEnv("DYNAMODB_REGION", "us_east-1"),
		},
		SNS: SNSConfig{
			Mock:     getEnv("SNS_MOCK", "true") == "true",
			TopicARN: getEnv("SNS_TOPIC_ARN", ""),
			Region:   getEnv("SNS_REGION", "us-east-1"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
