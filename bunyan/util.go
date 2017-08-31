package bunyan

import (
	"time"
	"fmt"
	"reflect"
	"regexp"
)

func StringDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func IntDefault(value int, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func NowTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

func TypeName(value interface{}) string {
	return fmt.Sprintf("%v", reflect.Type(value))
}

func IsHashMap(value interface{}) bool {
	r := regexp.MustCompile(`^*?map\[string\]`)
	return r.MatchString(TypeName(value))
}

func IsError(value interface{}) bool {
	r := regexp.MustCompile(`^*?errors.errorString`)
	return r.MatchString(TypeName(value))
}