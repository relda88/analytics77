package configuration

// Configuration holds the main application configuration mapping.
type Configuration struct {
	Server ConfigServer `json:"server"`
	Debug  ConfigDebug  `json:"debug"`
}

// ConfigServer stores the configuration for the primary application server.
type ConfigServer struct {
	Port int `json:"port"`
}

// ConfigDebug stores the configuration for the debugging interface.
type ConfigDebug struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}
