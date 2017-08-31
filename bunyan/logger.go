package bunyan

type Logger struct {
	config Config
	hostname string
	streams []Stream
}

// dynamically add a stream
func (l *Logger) AddStream(stream Stream) error {
	if err := stream.init(l.config); err != nil {
		return err
	}
	l.streams = append(l.streams, stream)
	return nil
}

// logging methods
func (l *Logger) Fatal(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelFatal {
			log.write(stream, LogLevelFatal)
		}
	}
}

func (l *Logger) Error(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelError {
			log.write(stream, LogLevelError)
		}
	}
}

func (l *Logger) Warn(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelWarn {
			log.write(stream, LogLevelWarn)
		}
	}
}

func (l *Logger) Info(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelInfo {
			log.write(stream, LogLevelInfo)
		}
	}
}

func (l *Logger) Debug(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelDebug {
			log.write(stream, LogLevelDebug)
		}
	}
}

func (l *Logger) Trace(args ...interface{}) {
	log := Log{ args: args, hostname: l.hostname }

	for _, stream := range l.streams {
		if stream.Level <= LogLevelTrace {
			log.write(stream, LogLevelTrace)
		}
	}
}