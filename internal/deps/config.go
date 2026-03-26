package deps

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	KafkaBrokers          string
	KafkaTopic            string
	KafkaGroupID          string
	KafkaUsername         string
	KafkaPassword         string
	KafkaSASLMechanism    string
	KafkaSecurityProtocol string

	HTTPPort string
	GRPCPort string
	LogLevel string

	PrometheusEnabled bool
	LokiEnabled       bool
	LokiURL           string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file .env not found, using environment variables")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	}

	return &Config{
		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBName:     v.GetString("DB_NAME"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),

		RedisHost:     v.GetString("REDIS_HOST"),
		RedisPort:     v.GetString("REDIS_PORT"),
		RedisPassword: v.GetString("REDIS_PASSWORD"),
		RedisDB:       v.GetInt("REDIS_DB"),

		KafkaBrokers:          v.GetString("KAFKA_BROKERS"),
		KafkaTopic:            v.GetString("KAFKA_TOPIC"),
		KafkaGroupID:          v.GetString("KAFKA_GROUP_ID"),
		KafkaUsername:         v.GetString("KAFKA_SASL_USERNAME"),
		KafkaPassword:         v.GetString("KAFKA_SASL_PASSWORD"),
		KafkaSASLMechanism:    v.GetString("KAFKA_SASL_MECHANISM"),
		KafkaSecurityProtocol: v.GetString("KAFKA_SECURITY_PROTOCOL"),

		HTTPPort: v.GetString("HTTP_PORT"),
		GRPCPort: v.GetString("GRPC_PORT"),
		LogLevel: strings.ToLower(v.GetString("LOG_LEVEL")),

		PrometheusEnabled: v.GetBool("PROMETHEUS_ENABLED"),
		LokiEnabled:       v.GetBool("LOKI_ENABLED"),
		LokiURL:           v.GetString("LOKI_URL"),
	}, nil
}
