package main

import (
	"errors"
	"github.com/bhoriuchi/go-bunyan/bunyan"
	"log"
	"os"
	"path/filepath"
)

/*
 * this file is only for quick testing of methods, and is not part of the actual project
 */
func main() {
	staticFields := make(map[string]interface{})
	staticFields["foo"] = "bar"
	baseDir := ""

	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		baseDir = wd
	}

	config := bunyan.Config{
		Name: "app",
		Streams: []bunyan.Stream{
			{
				Stream: os.Stdout,
			},
			{
				Path:  filepath.Join(baseDir, "info.log"),
				Level: bunyan.LogLevelInfo,
			},
			{
				Path:  filepath.Join(baseDir, "error.log"),
				Level: bunyan.LogLevelError,
			},
		},
		StaticFields: staticFields,
	}

	if l, err := bunyan.CreateLogger(config); err != nil {
		log.Fatal(err)
	} else {
		chStatic := make(map[string]interface{})
		chStatic["baz"] = "stuff"
		c := l.Child(chStatic)

		l.Info("i am the parent %d", 100)
		c.Info("i am the child")
		l.Error(errors.New("oh no!!!"), "what happened %s", "???")
	}
}
