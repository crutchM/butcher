package config

type Sentry struct {
	Enabled     bool   `yaml:"enabled"`
	DSN         string `yaml:"dsn"`
	Debug       bool   `yaml:"debug"`
	Environment string `yaml:"environment"`
}
