package deps

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/oauth"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

func KafkaClientOptions(cfg *Config, additionalOpts ...kgo.Opt) ([]kgo.Opt, error) {
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
	}

	opts = append(opts, additionalOpts...)

	securityProtocol := strings.ToUpper(cfg.KafkaSecurityProtocol)

	needsTLS := securityProtocol == "SSL" || securityProtocol == "SASL_SSL"
	needsSASL := strings.HasPrefix(securityProtocol, "SASL_")

	if needsTLS {
		tlsConfig, err := buildTLSConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build TLS config: %w", err)
		}
		opts = append(opts, kgo.DialTLSConfig(tlsConfig))
	} else if securityProtocol == "PLAINTEXT" || securityProtocol == "SASL_PLAINTEXT" {
		opts = append(opts, kgo.DialTLSConfig(nil))
	}

	if needsSASL && cfg.KafkaUsername != "" && cfg.KafkaPassword != "" {
		mechanism, err := buildSASLMechanism(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build SASL mechanism: %w", err)
		}
		opts = append(opts, kgo.SASL(mechanism))
	}

	return opts, nil
}

func buildTLSConfig(cfg *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cfg.KafkaTLSInsecureSkip,
	}

	if cfg.KafkaTLSCAFile != "" {
		caCert, err := os.ReadFile(cfg.KafkaTLSCAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	if cfg.KafkaTLSClientCert != "" && cfg.KafkaTLSClientKey != "" {
		cert, err := tls.LoadX509KeyPair(cfg.KafkaTLSClientCert, cfg.KafkaTLSClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

func buildSASLMechanism(cfg *Config) (sasl.Mechanism, error) {
	mechanism := strings.ToUpper(cfg.KafkaSASLMechanism)
	if mechanism == "" {
		mechanism = "PLAIN"
	}

	switch mechanism {
	case "PLAIN":
		return plain.Auth{
			User: cfg.KafkaUsername,
			Pass: cfg.KafkaPassword,
		}.AsMechanism(), nil

	case "SCRAM-SHA-256":
		return scram.Auth{
			User: cfg.KafkaUsername,
			Pass: cfg.KafkaPassword,
		}.AsSha256Mechanism(), nil

	case "SCRAM-SHA-512":
		return scram.Auth{
			User: cfg.KafkaUsername,
			Pass: cfg.KafkaPassword,
		}.AsSha512Mechanism(), nil

	case "OAUTHBEARER":
		return oauth.Auth{
			Token: cfg.KafkaPassword,
		}.AsMechanism(), nil

	case "AWS_MSK_IAM":
		return nil, fmt.Errorf("AWS MSK IAM mechanism not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported SASL mechanism: %s", mechanism)
	}
}
