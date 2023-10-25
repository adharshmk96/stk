package gsk

import "github.com/adharshmk96/stk/pkg/logging"

const (
	DEFAULT_PORT        = "8080"
	DEFAULT_STATIC_PATH = "/static"
	DEFAULT_STATIC_DIR  = "public/assets"
)

var DEFAULT_TEMPLATE_VARIABLES = map[string]interface{}{
	"Static": DEFAULT_STATIC_PATH,
}

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

	if initConfig.StaticPath == "" {
		initConfig.StaticPath = DEFAULT_STATIC_PATH
	}

	if initConfig.StaticDir == "" {
		initConfig.StaticDir = DEFAULT_STATIC_DIR
	}

	return initConfig
}
