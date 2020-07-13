package model

import (
	"encoding/json"

	"github.com/agalue/producer-enhancer/protobuf/producer"
)

// EnhancedEvent an object that represents an enhanced version of an OpenNMS Event
type EnhancedEvent struct {
	Event *producer.Event `json:"event,omitempty"`
	Node  *producer.Node  `json:"node,omitempty"`
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
