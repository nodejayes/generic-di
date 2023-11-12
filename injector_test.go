package di_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	di "github.com/nodejayes/generic-di"
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
		id     string
	}
)

func newConfiguration() *configuration {
	return &configuration{}
}

func newTextService() *textService {
	return &textService{
		config: di.Inject[configuration](),
		id:     uuid.NewString(),
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

func (ctx *textService) GetID() string {
	return ctx.id
}

func (ctx *messageService) GetTextServiceID() string {
	return ctx.texts.GetID()
}

func TestInject(t *testing.T) {
	msg := newMessageService()
	println(msg.texts.Greeting())
	if msg.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
}

func TestInject_Duplicate(t *testing.T) {
	msg1 := newMessageService()
	msg2 := newMessageService()
	println(msg1.texts.Greeting())
	println(msg2.texts.Greeting())
	if msg1.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if msg2.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if msg1.GetTextServiceID() != msg2.GetTextServiceID() {
		t.Errorf("expect same instance of textService")
	}
}

func TestInject_Parallel(t *testing.T) {
	for i := 0; i < 20; i++ {
		go func() {
			println(di.Inject[textService]().GetID())
		}()
	}
}

func TestInject_MultipleInstances(t *testing.T) {
	textServiceA := di.Inject[textService]("a")
	textServiceB := di.Inject[textService]("b")
	if textServiceA.GetID() == textServiceB.GetID() {
		t.Errorf("expect a seperate instance textServiceA and textServiceB but there was identical")
	}
}

func TestDestroy(t *testing.T) {
	_ = di.Inject[textService]("a")
	di.Destroy[textService]("a")
}
