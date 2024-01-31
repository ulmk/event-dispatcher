package event

import (
	"fmt"
	"testing"
)

type MyCommand struct{}

func (c MyCommand) Execute(params ...interface{}) error {
	for _, param := range params {
		fmt.Println(param)
	}
	return nil
}

func Test_Start(t *testing.T) {

	eventDispatcher := NewEventDispatcher()

	eventDispatcher.Sub("myTopic", MyCommand{})

	eventDispatcher.Pub("myTopic", "param1", 42, true)
}

type MyCommand2 struct {
	message string
}

func (c MyCommand2) Execute(params ...interface{}) error {
	if len(params) > 0 {
		fmt.Printf("%s: %v\n", c.message, params)
	} else {
		fmt.Println(c.message)
	}
	return nil
}

func Test_Command(t *testing.T) {

	eventDispatcher := NewEventDispatcher()

	err := eventDispatcher.Sub("example", MyCommand2{"Hello, World!"})
	if err != nil {
		fmt.Println("Subscription error:", err)
		return
	}

	// Publish event with parameters
	err = eventDispatcher.Pub("example", 42, "param")
	if err != nil {
		fmt.Println("Publish error:", err)
		return
	}
}
