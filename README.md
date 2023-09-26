# generic-di
Go Dependency Injection with Generics

## Example

configuration.go
```go
package main

import "github.com/nodejayes/generic-di"

func init() {
	// register the Struct Constructor Function for DI
	di.Injectable(NewConfiguration)
}

type Configuration struct {
	UserName string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		UserName: "Markus",
	}
}
```

greeter.go
```go
package main

import (
	"fmt"
	"github.com/nodejayes/generic-di"
)

func init() {
	di.Injectable(NewGreeter)
}

type Greeter struct {
	config *Configuration
}

func NewGreeter() *Greeter {
	return &Greeter{
		// here was the Configuration from configuration.go injected
		config: di.Inject[Configuration](),
	}
}

func (ctx *Greeter) Greet() string {
	return fmt.Sprintf("Hello, %s", ctx.config.UserName)
}
```

message_service.go
```go
package main

import "github.com/nodejayes/generic-di"

func init() {
	di.Injectable(NewMessageService)
}

type MessageService struct {
	greeter *Greeter
}

func NewMessageService() *MessageService {
	return &MessageService{
		// here was the Greeter from greeter.go injected
		greeter: di.Inject[Greeter](),
	}
}

func (ctx *MessageService) Welcome() string {
	return ctx.greeter.Greet()
}
```

main.go

```go
package main

import di "github.com/nodejayes/generic-di"

func main() {
	msgService := di.Inject[MessageService]()
	// prints the message "Hello, Markus"
	println(msgService.Welcome())
}
```
