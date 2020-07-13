package model

import (
	"encoding/json"
	"log"

	"github.com/agalue/producer-enhancer/protobuf/producer"
	"github.com/lovoo/goka"
)

// EnhancedEvent an object that represents an enhanced version of an OpenNMS Event
type EnhancedEvent struct {
	Event *producer.Event `json:"event,omitempty"`
	Node  *producer.Node  `json:"node,omitempty"`
}

func (e EnhancedEvent) enhance(event *producer.Event, ctx goka.Context, nodesTopic goka.Table) *EnhancedEvent {
	enhanced := &EnhancedEvent{Event: event}
	log.Printf("enhancing event with ID %d and UEI %s", event.Id, event.Uei)
	if c := event.NodeCriteria; c != nil {
		node := ctx.Lookup(nodesTopic, producer.GetNodeKey(c))
		if node != nil {
			enhanced.Node = node.(*producer.Node)
		}
	}
	return enhanced
}

// DefineGroup creates a Goke GroupGraph to enhance events
func (e EnhancedEvent) DefineGroup(groupID, eventsTopic, nodesTopic, targetTopic string) *goka.GroupGraph {
	return goka.DefineGroup(goka.Group(groupID),
		goka.Input(goka.Stream(eventsTopic), new(producer.EventCodec), func(ctx goka.Context, msg interface{}) {
			event := msg.(*producer.Event)
			enhanced := e.enhance(event, ctx, goka.Table(nodesTopic))
			ctx.Emit(goka.Stream(targetTopic), ctx.Key(), enhanced)
		}),
		goka.Output(goka.Stream(targetTopic), new(EnhancedEventCodec)),
		goka.Lookup(goka.Table(nodesTopic), new(producer.NodeCodec)),
	)
}

// EnhancedEventCodec a Goka Codec to encode/decode an enhanced Event to/from JSON
type EnhancedEventCodec struct{}

// Encode encodes a Protobuf byte array into an OpenNMS Event object
func (c EnhancedEventCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = json.Marshal(value)
	return data, err
}

// Decode decodes a Protobuf byte array into an OpenNMS Event object
func (c EnhancedEventCodec) Decode(data []byte) (value interface{}, err error) {
	value = &EnhancedEvent{}
	err = json.Unmarshal(data, value.(*EnhancedEvent))
	return value, err
}
