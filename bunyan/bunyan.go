package bunyan

import (
	"fmt"
	"regexp"
	"os"
	"io"
)

// export the log version
const LOG_VERSION = 0

type Config struct {
	Name string
	Level string
	Stream io.Writer
	Streams []Stream
	Serializers map[string]func(value interface{}) interface{}
	StaticFields map[string]interface{}
}

// main function to create a new logger
func CreateLogger(args ...interface{}) (bunyanLogger, error) {
	config := Config{}
	logger := bunyanLogger{}
	r := regexp.MustCompile(`bunyan.Config$`)

	if len(args) == 0 {
		return logger, fmt.Errorf("Create logger requires either a bunyan.Config or String argument")
	}

	// get hostname
	if hostname, err := os.Hostname(); err != nil {
		return logger, err
	} else {
		logger.hostname = hostname
	}

	arg := args[0]
	argType := typeName(arg)

	if argType == "string" {
		config.Name = arg.(string)
	} else if r.MatchString(argType) {
		config = arg.(Config)
		if config.Name == "" {
			return logger, fmt.Errorf("Bunyan Config requires a name, none specified")
		}
	} else {
		return logger, fmt.Errorf("Create logger requires either a bunyan.Config or String argument")
	}

	// add the config to the logger
	logger.config = config
	logger.staticFields = config.StaticFields
	logger.serializers = config.Serializers

	// add the streams
	if len(config.Streams) != 0 {
		for _, stream := range config.Streams {
			logger.AddStream(stream)
		}
	} else if config.Stream != nil {
		simpleStream := Stream{ Stream: config.Stream, Name: config.Name }
		logger.AddStream(simpleStream)
	}

	return logger, nil
}