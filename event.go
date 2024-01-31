package event

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrCommandOrTopicNil = errors.New("command or topic couldn't be empty")
	ErrTopicNotFound     = errors.New("topic not found")
)

type Command interface {
	Execute(params ...interface{}) error
}

type Event struct {
	Cmd   Command
	Topic string
}

type EventDispatcher struct {
	mu     sync.RWMutex
	events map[string][]Event
}

// var eventDispatcher = &EventDispatcher{
// 	events: make(map[string][]Event),
// }

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		events: make(map[string][]Event),
	}
}

func (e *EventDispatcher) Sub(topic string, command Command) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if command == nil || topic == "" {
		return ErrCommandOrTopicNil
	}

	event := Event{
		Topic: topic,
		Cmd:   command,
	}

	e.events[topic] = append(e.events[topic], event)

	return nil
}

func (e *EventDispatcher) Unsub(topic string, command Command) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if topic == "" {
		return ErrCommandOrTopicNil
	}
	events, found := e.events[topic]
	if !found {
		return ErrTopicNotFound
	}

	var updatedEvents []Event

	for _, event := range events {
		if event.Cmd != command {
			updatedEvents = append(updatedEvents, event)
		}
	}

	if len(updatedEvents) == len(events) {
		return ErrTopicNotFound
	}

	e.events[topic] = updatedEvents

	return nil
}

func (e *EventDispatcher) Pub(topic string, params ...interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	events := e.events[topic]

	if len(events) == 0 {
		return ErrTopicNotFound
	}
	var wg sync.WaitGroup
	// for i := len(commands) - 1; i >= 0; i-- {
	// 	go commands[i].Cmd.Execute(params...)
	// }
	for _, event := range events {
		wg.Add(1)
		go func(e Event) {
			defer wg.Done()
			err := e.Cmd.Execute(params...)
			if err != nil {
				fmt.Printf("Error executing command for topic %s: %v\n", e.Topic, err)
			}
		}(event)
	}

	wg.Wait()

	return nil
}
