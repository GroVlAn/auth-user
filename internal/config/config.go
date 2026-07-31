package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type HTTP struct {
	Port              string        `yaml:"port" env-default:"8080"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes" env-default:"4096"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env-default:"10s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env-default:"10s"`
	BaseHTTPPath      string        `yaml:"base_http_path" env-default:"/api"`
}

type Middleware struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type GRPC struct {
	Port          string `yaml:"port"`
	AccessApiHost string `yaml:"access_api_host"`
	AccessApiPort string `yaml:"access_api_port"`
}

type Settings struct {
	DefaultTimeout time.Duration `yaml:"default_timeout"`
	HashCost       int           `yaml:"hash_cost"`
}

type PostgresSettings struct {
	Host    string `yaml:"host" env-required:"true"`
	Port    string `yaml:"port"`
	SSLMode string `yaml:"ssl_mode"`
}

type Cache struct {
	DefaultExpiration time.Duration `yaml:"default_expiration"`
	CleanupInterval   time.Duration `yaml:"cleanup_interval"`
	UserTTL           time.Duration `yaml:"user_ttl"`
	RoleTTL           time.Duration `yaml:"role_ttl"`
}

type Vault struct {
	SecretToken string `env:"VAULT_SECRET_TOKEN" env-required:"true"`
	Address     string `env:"VAULT_ADDRESS" env-required:"true"`
	Mount       string `env:"VAULT_MOUNT" env-required:"true"`
}

type VaultPaths struct {
	Postgres string `env:"POSTGRES_PATH" env-required:"true"`
	Hasher   string `env:"HASHER_PATH" env-required:"true"`
}

type Config struct {
	HTTP       HTTP             `yaml:"http"`
	Middleware Middleware       `yaml:"middleware"`
	GRPC       GRPC             `yaml:"grpc"`
	DB         PostgresSettings `yaml:"db"`
	Settings   Settings         `yaml:"settings"`
	Cache      Cache            `yaml:"cache"`
	Vault      Vault
	VaultPaths VaultPaths
}

func New(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	return cfg, nil
}

func LoadEnv(filenames ...string) error {
	if len(filenames) == 0 {
		return godotenv.Load()
	}

	for _, filename := range filenames {
		if err := godotenv.Load(filename); err != nil {
			return fmt.Errorf("loading env file %s: %w", filename, err)
		}
	}
	return nil
}
