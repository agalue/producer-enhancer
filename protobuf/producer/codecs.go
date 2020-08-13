package producer

import (
	"github.com/golang/protobuf/proto"
)

// EventCodec a Goka Codec to encode/decode an OpenNMS Event to/from Protobuf
type EventCodec struct{}

// Encode encodes an OpenNMS Event object into a Protobuf byte array
func (c EventCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = proto.Marshal(value.(*Event))
	return data, err
}

// Decode decodes a Protobuf byte array into an OpenNMS Event object
func (c EventCodec) Decode(data []byte) (value interface{}, err error) {
	value = &Event{}
	err = proto.Unmarshal(data, value.(*Event))
	return value, err
}

// AlarmCodec a Goka Codec to encode/decode an OpenNMS Alarm to/from Protobuf
type AlarmCodec struct{}

// Encode encodes an OpenNMS Alarm object into a Protobuf byte array
func (c AlarmCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = proto.Marshal(value.(*Alarm))
	return data, err
}

// Decode decodes a Protobuf byte array into an OpenNMS Alarm object
func (c AlarmCodec) Decode(data []byte) (value interface{}, err error) {
	value = &Alarm{}
	err = proto.Unmarshal(data, value.(*Alarm))
	return value, err
}

// NodeCodec a Goka Codec to encode/decode an OpenNMS Node to/from Protobuf
type NodeCodec struct{}

// Encode encodes an OpenNMS Node object into a Protobuf byte array
func (c NodeCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = proto.Marshal(value.(*Node))
	return data, err
}

// Decode decodes a Protobuf byte array into an OpenNMS Node object
func (c NodeCodec) Decode(data []byte) (value interface{}, err error) {
	value = &Node{}
	err = proto.Unmarshal(data, value.(*Node))
	return value, err
}
