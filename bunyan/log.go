package bunyan

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	//"reflect"
)

type bunyanLog struct {
	args   []interface{}
	logger Logger
}

// serializes a log field
func (l *bunyanLog) serialize(key string, value interface{}) interface{} {
	if fn, ok := l.logger.serializers[key]; ok {
		return fn(value)
	} else if isError(value) {
		return fmt.Sprintf("%v", value)
	} else {
		return value
	}
}

// prints a formatted string using the arguments provided
func (l *bunyanLog) sprintf(args []interface{}) string {
	if typeName(args[0]) != "string" {
		return fmt.Sprintf("", args[0:]...)
	}
	return fmt.Sprintf(args[0].(string), args[1:]...)
}

// composes a json formatted log string and writes it to the appropriate stream
func (l *bunyanLog) write(stream Stream, level int) error {
	data := make([]byte, 0)
	argl := len(l.args)

	if argl == 0 {
		return nil
	}
	d := make(map[string]interface{})
	defaultKey:=l.logger.defaultKey
	d["v"] = 0
	d["level"] = level
	d["name"] = stream.Name
	d["hostname"] = l.logger.hostname
	d["pid"] = os.Getppid()
	d["time"] = nowTimestamp()
	d["details"]=make(map[string]map[string]string)
	details:=make(map[string]map[string]interface{})
	// add static fields first
	for key, value := range l.logger.staticFields {
		if canSetField(key) {
			d[key] = l.serialize(key, value)
		}
	}


	// add passed fields/data last
	for _, element := range l.args {

		if typeName(element) == "string" {
			// if  argument that is a string, the string is the msg.
			//returning concat string if string arguments more than one
			var newMsg string
			if d[defaultKey] !=nil {
				newMsg = d[defaultKey].(string) +", " + element.(string)	
			}else{
				newMsg = element.(string)
			}			
			d[defaultKey] = l.serialize(defaultKey, newMsg)
		} else if isError(element) {
			// if the  argument is an error, set error field with string value of error
			d["error"] = l.serialize("error", element)
		} else if isHashMap(element) {
			// if the  argument is a hashmap, process its values
			for key, value := range element.(map[string]interface{}) {
				if canSetField(key) {
					d[key] = l.serialize(key, value)
				}
			}
		} else if isStruct(element) {
			// get details with reflect value if the argumant is struct.
			// If the struct is more than one, it returns the value with the struct name.
			detail := getDetailsLog(element)			
			for key,value:=range detail{		
				details[key]=make(map[string]interface{})
				details[key]=value.(map[string]interface{})
			}
		} else {
			d[defaultKey] = l.serialize(defaultKey, l.sprintf([]interface{}{d[defaultKey],element}))
		}
	}

	d["details"]=details
	// marshal the json
	if jsonData, err := json.Marshal(d); err != nil {
		return err
	} else {
		data = []byte(fmt.Sprintf("%s\n", string(jsonData)))
	}

	switch stream.Type {
	case LogTypeStream:
		return l.writeStream(stream, data)
	case LogTypeFile:
		return l.writeFile(stream, data)
	case LogTypeRotatingFile:
		return l.writeRotatingFile(stream, data)
	case LogTypeRaw:
		return l.writeStream(stream, data)
	}
	return nil
}

// writes the data to a stream that implements io.Writer
func (l *bunyanLog) writeStream(stream Stream, data []byte) error {
	stream.Stream.Write(data)
	return nil
}

// writes the data to a log file
func (l *bunyanLog) writeFile(stream Stream, data []byte) error {
	if f, err := os.OpenFile(stream.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Printf("[bunyan] error: %v", err)
	} else if _, err := f.Write(data); err != nil {
		log.Printf("[bunyan] error: %v", err)
	} else if err := f.Close(); err != nil {
		log.Printf("[bunyan] error: %v", err)
	}
	return nil
}

// TODO: implement this
func (l *bunyanLog) writeRotatingFile(stream Stream, data []byte) error {
	return nil
}
