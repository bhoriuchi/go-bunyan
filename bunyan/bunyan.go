package bunyan

import (
	"fmt"
	"regexp"
	"os"
)

func createLogger(args ...interface{}) (Logger, error) {
	config := Config{}
	logger := Logger{}
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
	argType := TypeName(arg)

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

	// add the streams
	if len(config.Streams) != 0 {
		logger.streams = config.Streams
	} else if config.Stream != nil {
		logger.streams = append(logger.streams, Stream{ Stream: config.Stream, Name: config.Name })
	}

	// add the config to the logger
	logger.config = config

	// init all the streams
	for _, stream := range logger.streams {
		if err := stream.init(config); err != nil {
			return logger, err
		}
	}

	return logger, nil
}