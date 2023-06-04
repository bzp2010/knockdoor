package config

// Config is the entrypoint for all configuration items
type Config struct {
	Log  Log  `mapstructure:"log"`
	Port Port `mapstructure:"port"`
}

// Log configuration
type Log struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// Port is knockdoor core configuration
type Port struct {
	Mode         string    `mapstructure:"mode"`
	Static       *[]uint16 `mapstructure:"static"`
	TOTP         *PortTOTP `mapstructure:"totp"`
	SkipLoopback bool      `mapstructure:"skipLoopback"`
}

// PortTOTP is knockdoor core configuration for TOTP mode
type PortTOTP struct {
	Secret string `mapstructure:"secret"`
	Prefix string `mapstructure:"prefix"`
}
