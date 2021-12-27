package bunyan

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func stringDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func intDefault(value int, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func nowTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

func typeName(value interface{}) string {
	return fmt.Sprintf("%T", value)
}

func isHashMap(value interface{}) bool {
	r := regexp.MustCompile(`^*?map\[string\]`)
	return r.MatchString(typeName(value))
}

func isStruct(value interface{}) bool {
	return (reflect.ValueOf(value).Kind() == reflect.Struct) || (reflect.ValueOf(value).Kind() == reflect.Ptr)
}

func getDetailsLog(value interface{}) map[string]interface{} {

	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	details := make(map[string]interface{})
	if v.Kind() == reflect.Struct {
		structS := make(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			structS[v.Type().Field(i).Name] = getValue(v.Field((i)))
		}
		details[v.Type().Name()] = structS
		return details
	} else {
		details[v.Type().Name()] = getValue(v)
	}
	return details
}

func getValue(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Bool:
		return value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Slice:
		s := value
		var c []interface{}
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Kind() == reflect.Struct {
				existInterface := value.CanInterface()
				if existInterface {
					c = append(c, getDetailsLog(s.Index(i).Interface()))
				} else {
					c = append(c, getDetailsLog(s.Index(i)))
				}
			} else if s.Index(i).Kind() == reflect.Ptr {
				existInterface := value.CanInterface()
				if existInterface {
					c = append(c, getDetailsLog(s.Index(i).Interface()))
				} else {
					c = append(c, getDetailsLog(s.Index(i)))
				}
			} else {
				c = append(c, getValue(s.Index(i)))
			}
		}
		return c
	case reflect.Ptr:
		t := value.CanInterface()
		if t {
			return getDetailsLog(value.Interface())
		} else {
			if value.Kind() == reflect.Ptr {
				value = value.Elem()
			}
			var address string
			if value.CanAddr() {
				address = fmt.Sprint(value.Addr())
			}
			return address
		}
	case reflect.String:
		return value.String()
	case reflect.Struct:
		existInterface := value.CanInterface()
		if existInterface {
			interfacesV := value.Interface()
			return getDetailsLog(interfacesV)
		} else {
			return getDetailsLog(reflect.ValueOf(value))
		}
	default:
		return value
	}
}

func isError(value interface{}) bool {
	r := regexp.MustCompile(`^*?errors.errorString`)
	return r.MatchString(typeName(value))
}

func canSetField(key interface{}) bool {
	switch strings.ToLower(key.(string)) {
	case "v":
		return false
	case "level":
		return false
	case "details":
		return false
	default:
		return true
	}
}

func toLogLevelInt(level string) int {
	switch strings.ToLower(level) {
	case LogLevelFatal:
		return 60
	case LogLevelError:
		return 50
	case LogLevelWarn:
		return 40
	case LogLevelInfo:
		return 30
	case LogLevelDebug:
		return 20
	case LogLevelTrace:
		return 10
	default:
		return 0
	}
}
