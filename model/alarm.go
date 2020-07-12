package model

import (
	"encoding/json"

	"github.com/agalue/producer-enhancer/protobuf/producer"
)

// EnhancedAlarm an object that represents an enhanced version of an OpenNMS Alarm
type EnhancedAlarm struct {
	Alarm *producer.Alarm `json:"alarm,omitempty"`
	Node  *producer.Node  `json:"node,omitempty"`
}

// EnhancedAlarmCodec a Goka Codec to encode/decode an enhanced alarm to/from JSON
type EnhancedAlarmCodec struct{}

// Encode encodes a Protobuf byte array into an OpenNMS Alarm object
func (c EnhancedAlarmCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = json.Marshal(value)
	return data, err
}

// Decode decodes a Protobuf byte array into an OpenNMS Alarm object
func (c EnhancedAlarmCodec) Decode(data []byte) (value interface{}, err error) {
	value = &EnhancedAlarm{}
	err = json.Unmarshal(data, value.(*EnhancedAlarm))
	return value, err
}
