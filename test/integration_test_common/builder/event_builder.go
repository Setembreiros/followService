package integration_test_builder

import (
	"encoding/json"
	"followservice/internal/bus"
	"testing"

	"github.com/stretchr/testify/assert"
)

type EventBuilder struct {
	t    *testing.T
	name string
	data []byte
}

func NewEventBuilder(t *testing.T) *EventBuilder {
	return &EventBuilder{
		t:    t,
		name: "",
		data: []byte{},
	}
}

func (eb *EventBuilder) WithName(name string) *EventBuilder {
	eb.name = name
	return eb
}

func (eb *EventBuilder) WithData(data any) *EventBuilder {
	dataEvent, err := serializeData(data)
	assert.Nil(eb.t, err)
	eb.data = dataEvent

	return eb
}

func (eb *EventBuilder) Build() *bus.Event {
	return &bus.Event{
		Type: eb.name,
		Data: eb.data,
	}
}

func serializeData(data any) ([]byte, error) {
	return json.Marshal(data)
}
