package di_test

import (
	"fmt"
	"github.com/nodejayes/webtools/di"
	"testing"
)

func init() {
	di.Injectable(newTextService)
	di.Injectable(newMessageService)
	di.Injectable(newConfiguration)
}

type (
	configuration  struct{}
	messageService struct {
		texts *textService
	}

	textService struct {
		config *configuration
	}
)

func newConfiguration() *configuration {
	return &configuration{}
}

func newTextService() *textService {
	return &textService{
		config: di.Inject[configuration](),
	}
}

func newMessageService() *messageService {
	return &messageService{
		texts: di.Inject[textService](),
	}
}

func (ctx *configuration) GetUserName() string {
	return "Markus"
}

func (ctx *textService) Greeting() string {
	return fmt.Sprintf("Hello %s", ctx.config.GetUserName())
}

func TestInject(t *testing.T) {
	msg := newMessageService()
	println(msg.texts.Greeting())
	if msg.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
}
