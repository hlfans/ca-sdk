package config

type CAConfig struct {
	Host string    `yaml:"host"`
	Tls  TlsConfig `yaml:"tls"`
}
