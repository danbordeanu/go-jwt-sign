package configuration

import "github.com/danbordeanu/go-utils"

type Configuration struct {
	Swagger CSwagger

	IngressHost   string
	IngressPrefix string

	// Dependencies
	JaegerEndpoint string

	// jaeger
	JaegerEngine string

	// Configuration
	HttpPort int32

	// Internal settings
	CleanupTimeoutSec int32
	Environment       string
	UseTelemetry      string
	Development       bool
	GinLogger         bool
	UseSwagger        bool
	VaultLogging      bool
	Initialized       bool

	// baseUrl page
	RequestBaseUrl string

	// Cors allow origins
	CorsAllowOrigins string
}

var appConfig Configuration

func AppConfig() *Configuration {
	if appConfig.Initialized == false {
		loadEnvironmentVariables()
		appConfig.Initialized = true
	}
	return &appConfig
}

// loadEnvironmentVariables load env variables
func loadEnvironmentVariables() {
	// jaeger telemetry settings
	appConfig.JaegerEngine = utils.EnvOrDefault("JAEGER_ENGINE_NAME", "http://localhost:14268/api/traces")
	appConfig.Environment = utils.EnvOrDefault("ENVIRONMENT", "local")
	appConfig.JaegerEndpoint = utils.EnvOrDefault("JAEGER_ENDPOINT", "")
	appConfig.CleanupTimeoutSec = utils.EnvOrDefaultInt32("SHUTDOWN_TIMEOUT", 300)
	appConfig.IngressHost = utils.EnvOrDefault("INGRESS_HOST", "jwt-sign")
	appConfig.IngressPrefix = utils.EnvOrDefault("INGRESS_PREFIX", "")
	// request base url
	appConfig.RequestBaseUrl = utils.EnvOrDefault("REQUEST_BASE_URL", "http://localhost:8080")

	// CORS allow origins
	appConfig.CorsAllowOrigins = utils.EnvOrDefault("CORS_ALLOW_ORIGINS", "Disabled")

}
