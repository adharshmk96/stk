package gsk

import "github.com/adharshmk96/stk/pkg/logging"

// Initialize the server configurations
// if no configurations are passed, default values are used
func initConfig(config ...*ServerConfig) *ServerConfig {
	var initConfig *ServerConfig
	if len(config) == 0 {
		initConfig = &ServerConfig{}
	} else {
		initConfig = config[0]
	}

	if initConfig.Port == "" {
		initConfig.Port = "8080"
	}

	if initConfig.Logger == nil {
		initConfig.Logger = logging.NewSlogLogger()
	}

	if initConfig.BodySizeLimit == 0 {
		initConfig.BodySizeLimit = 1
	}

	return initConfig
}
