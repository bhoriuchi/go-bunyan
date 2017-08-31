package bunyan

import (
	"io"
	"fmt"
)

type Stream struct {
	// universal fields
	Type string
	Level string
	Name string

	// stream fields
	Stream io.Writer

	// file fields
	Path string

	// rotating file fields
	Period string
	Count int
}

func (s *Stream) init(config Config) error {
	if s.Type != LogTypeStream && s.Type != LogTypeFile && s.Type != LogTypeRotatingFile && s.Type != LogTypeRaw {
		if s.Stream != nil {
			s.Type = LogTypeStream
		} else if s.Path != "" {
			s.Type = LogTypeFile
		}
	}

	if s.Type == "" {
		return fmt.Errorf("Invalid stream options, could not determine stream type")
	}

	s.Name = stringDefault(s.Name, config.Name)
	s.Level = stringDefault(s.Level, config.Level)

	// set default log level
	if toLogLevelInt(s.Level) <= 0 {
		s.Level = LogLevelInfo
	}

	// check the rest of the types
	switch s.Type {
	case LogTypeStream:
		if s.Stream == nil {
			return fmt.Errorf("Stream logs require the %q argument implement interface %q", "Stream", "io.Writer")
		}
		break
	case LogTypeFile:
		if s.Path == "" {
			return fmt.Errorf("File logs require the %q argument", "Path")
		}
		break
	case LogTypeRotatingFile:
		if s.Path == "" {
			return fmt.Errorf("Rotating File logs require the %q argument", "Path")
		}
		if s.Period == "" {
			s.Period = "1d"
		}
		if s.Count == 0 {
			s.Count = 10
		}
		break
	case LogTypeRaw:
		if s.Stream == nil {
			return fmt.Errorf("Raw logs require the %q argument implement interface %q", "Stream", "io.Writer")
		}
	}

	return nil
}
