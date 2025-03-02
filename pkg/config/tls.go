package config

type TlsConfig struct {
	Enabled    bool `yaml:"enabled"`
	SkipVerify bool `yaml:"skip_verify"`

	// Cert take precedence over CertPath
	Cert     []byte `yaml:"cert"`
	CertPath string `yaml:"cert_path"`

	// Key take precedence over KeyPath
	Key     []byte `yaml:"key"`
	KeyPath string `yaml:"key_path"`

	// CACert take precedence over CACertPath
	CACert     []byte `yaml:"ca_cert"`
	CACertPath string `yaml:"ca_cert_path"`
}
