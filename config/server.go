package config

var ServerConfigValues ServerConfig

// Model that links to config.yml file
type ServerConfig struct {
	Minio struct {
		Endpoint string    `yaml:"endpoint" env:"ENDPOINT" env-description:"Minio Server Endpoint"`
		AccessKeyID        string `yaml:"access-key-id" env:"ACCESS_KEY_ID" env-description:"Access Key ID"`
		SecretAccessKey        string `yaml:"secret-access-key" env:"SECRET_ACCESS_KEY" env-description:"Secret Access Key"`
	} `yaml:"minio"`
	Server struct {
		ApiPath            string   `yaml:"api-path"  env:"API_PATH" env-description:"API base path"`
		ApiVersion         string   `yaml:"api-version"  env:"API_VERSION" env-description:"API Version"`
		CorsAllowedClients []string `yaml:"cors-allowed-clients" env:"CORS_ALLOWED_CLIENTS"  env-description:"List of allowed CORS Clients"`
		Environment        string   `yaml:"environment" env:"SERVER_ENVIRONMENT"  env-description:"server environment"`

		Host     string `yaml:"host"  env:"SERVER_HOST" env-description:"server host"`
		Port     string `yaml:"port" env:"SERVER_PORT"  env-description:"server port"`
		Protocol string `yaml:"protocol" env:"SERVER_PROTOCOL"  env-description:"server protocol"`
	} `yaml:"server"`
}
