# go-bunyan
Bunyan logger written in Golang. 
See https://github.com/trentm/node-bunyan for in depth documentation. 
This package attempts to implement the original node.js library as closely
as possible using golang.

## TODO
* Implement `rotating-file` stream
* Write tests
* Publish to go registry

### Simple usage

```golang
import (
    "os"
    "github.com/bhoriuchi/go-bunyan/bunyan"
)

func main() {
    config := bunyan.Config{
        Name: "app",
        Stream: os.Stdout,
        Level: bunyan.LogLevelDebug
    }
    
    if log, err := bunyan.CreateLogger(config); err == nil {
        log.Info("Hello %s!", "World")
    }
}
```

### Multi-stream usage

```golang
import (
    "os"
    "github.com/bhoriuchi/go-bunyan/bunyan"
)

func main() {
	staticFields := make(map[string]interface{})
	staticFields["foo"] = "bar"
	
    config := bunyan.Config{
        Name: "app",
        Streams: []bunyan.Stream{
            {
                Name: "app-info",
                Level: bunyan.LogLevelInfo,
                Stream: os.Stdout,
            },
            {
                Name: "app-errors",
                Level: bunyan.LogLevelError,
                Path: "/path/to/logs/app-errors.log"
            },
        },
        StaticFields: staticFields,
    }
    
    if log, err := bunyan.CreateLogger(config); err == nil {
        log.Info("Hello %s!", "World")
    }
}
```