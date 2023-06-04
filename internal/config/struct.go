package config

// Config is the entrypoint for all configuration items
type Config struct {
	Log   Log   `mapstructure:"log"`
	Knock Knock `mapstructure:"knock"`
	Door  Door  `mapstructure:"door"`
}

// Log configuration
type Log struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// Knock is knockdoor knock configuration
type Knock struct {
	Mode         string     `mapstructure:"mode"`
	Static       *[]uint16  `mapstructure:"static"`
	TOTP         *KnockTOTP `mapstructure:"totp"`
	SkipLoopback bool       `mapstructure:"skipLoopback"`
}

// KnockTOTP is knockdoor core configuration for TOTP mode
type KnockTOTP struct {
	Secret string `mapstructure:"secret"`
	Prefix string `mapstructure:"prefix"`
}

// Door is knockdoor door configuration
type Door struct {
	Type     string        `mapstructure:"type"`
	RouterOS *DoorRouterOS `mapstructure:"routeros"`
}

// DoorRouterOS is knockdoor door configuration for RouterOS
type DoorRouterOS struct {
	Endpoint        string `mapstructure:"endpoint"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	AddressListName string `mapstructure:"addressListName"`
}
