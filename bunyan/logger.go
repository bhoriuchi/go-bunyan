package bunyan

type bunyanLogger struct {
	config Config
	hostname string
	streams []Stream
	staticFields map[string]interface{}
	serializers map[string]func(value interface{}) interface{}
}

// dynamically add a stream
func (l *bunyanLogger) AddStream(stream Stream) error {
	if err := stream.init(l.config); err != nil {
		return err
	}
	l.streams = append(l.streams, stream)
	return nil
}

// dynamically add serializers
func (l *bunyanLogger) AddSerializers(serializers map[string]func(value interface{}) interface{}) {
	for key, value := range serializers {
		l.serializers[string(key)] = value
	}
}

// TODO: implement this
func (l *bunyanLogger) Level() {

}

// creates a new child logger with extra static fields
func (l *bunyanLogger) Child(fields map[string]interface{}) bunyanLogger {
	newStaticFields := make(map[string]interface{})

	// merge the static fields into the new logger
	for key, field := range l.staticFields {
		newStaticFields[string(key)] = field
	}
	for key, field := range fields {
		newStaticFields[string(key)] = field
	}

	logger := bunyanLogger{
		config: l.config,
		hostname: l.hostname,
		streams: l.streams,
		staticFields: newStaticFields,
	}
	return logger
}

// logging methods
func (l *bunyanLogger) Fatal(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelFatal) {
			_log.write(stream, toLogLevelInt(LogLevelFatal))
		}
	}
}

func (l *bunyanLogger) Error(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelError) {
			_log.write(stream, toLogLevelInt(LogLevelError))
		}
	}
}

func (l *bunyanLogger) Warn(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelWarn) {
			_log.write(stream, toLogLevelInt(LogLevelWarn))
		}
	}
}

func (l *bunyanLogger) Info(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelInfo) {
			_log.write(stream, toLogLevelInt(LogLevelInfo))
		}
	}
}

func (l *bunyanLogger) Debug(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelDebug) {
			_log.write(stream, toLogLevelInt(LogLevelDebug))
		}
	}
}

func (l *bunyanLogger) Trace(args ...interface{}) {
	_log := bunyanLog{ args: args, logger: *l }

	for _, stream := range l.streams {
		if toLogLevelInt(stream.Level) <= toLogLevelInt(LogLevelTrace) {
			_log.write(stream, toLogLevelInt(LogLevelTrace))
		}
	}
}