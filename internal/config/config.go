package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {
	config = viper.New()
	config.SetConfigFile(".env")
	config.SetConfigType("env")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Warning: failed to read .env file: %v", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}

// Database
func DBHost() string {
	return config.GetString("DB_HOST")
}

func DBPort() string {
	return config.GetString("DB_PORT")
}

func DBName() string {
	return config.GetString("DB_NAME")
}

func DBUser() string {
	return config.GetString("DB_USER")
}

func DBPassword() string {
	return config.GetString("DB_PASSWORD")
}

func DBSSLMode() string {
	return config.GetString("DB_SSLMODE")
}

// Redis
func RedisHost() string {
	return config.GetString("REDIS_HOST")
}

func RedisPort() string {
	return config.GetString("REDIS_PORT")
}

func RedisPassword() string {
	return config.GetString("REDIS_PASSWORD")
}

func RedisDB() int {
	return config.GetInt("REDIS_DB")
}

// Kafka
func KafkaBrokers() string {
	return config.GetString("KAFKA_BROKERS")
}

func KafkaTopicDocument() string {
	return config.GetString("KAFKA_TOPIC_DOCUMENT")
}

func KafkaGroupID() string {
	return config.GetString("KAFKA_GROUP_ID")
}

// Server
func HTTPPort() string {
	return config.GetString("HTTP_PORT")
}

func GRPCPort() string {
	return config.GetString("GRPC_PORT")
}

func LogLevel() string {
	return strings.ToLower(config.GetString("LOG_LEVEL"))
}

func Environment() string {
	return config.GetString("ENVIRONMENT")
}

func IsDevelopment() bool {
	return Environment() == "development"
}

func IsProduction() bool {
	return Environment() == "production"
}

// Monitoring
func PrometheusEnabled() bool {
	return config.GetBool("PROMETHEUS_ENABLED")
}

func LokiEnabled() bool {
	return config.GetBool("LOKI_ENABLED")
}

func LokiURL() string {
	return config.GetString("LOKI_URL")
}
