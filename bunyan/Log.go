package bunyan

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
)

type Log struct {
	args []interface{}
	hostname string
	data []byte
}

func (l *Log) write(stream Stream, level int) error {
	argl := len(l.args)
	if argl == 0 {
		return nil
	}

	d := make(map[string]interface{})
	d["v"] = 0
	d["level"] = level
	d["name"] = stream.Name
	d["hostname"] = l.hostname
	d["pid"] = os.Getppid()
	d["time"] = NowTimestamp()

	if argl == 1 && TypeName(l.args[0]) == "string" {
		d["msg"] = l.args[0].(string)
	}


	// marshal the json
	if jsonData, err := json.Marshal(d); err != nil {
		return err
	} else {
		l.data = []byte(fmt.Sprintf("%s\n", string(jsonData)))
	}

	switch stream.Type {
	case LogTypeStream:
		return l.writeStream(stream)
	case LogTypeFile:
		return l.writeFile(stream)
	case LogTypeRotatingFile:
		return l.writeRotatingFile(stream)
	case LogTypeRaw:
		return l.writeStream(stream)
	}
	return nil
}

func (l *Log) writeStream(stream Stream) error {
	stream.Stream.Write(l.data)
	return nil
}

func (l *Log) writeFile(stream Stream) error {
	if f, err := os.OpenFile(stream.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Printf("[bunyan] error: %v", err)
	} else if _, err := f.Write(l.data); err != nil {
		log.Printf("[bunyan] error: %v", err)
	} else if err := f.Close(); err != nil {
		log.Printf("[bunyan] error: %v", err)
	}
	return nil
}

func (l *Log) writeRotatingFile(stream Stream) error {
	return nil
}