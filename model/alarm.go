package model

import (
	"encoding/json"
	"log"

	"github.com/agalue/producer-enhancer/protobuf/producer"
	"github.com/lovoo/goka"
)

// EnhancedAlarm an object that represents an enhanced version of an OpenNMS Alarm
type EnhancedAlarm struct {
	Alarm *producer.Alarm `json:"alarm,omitempty"`
	Node  *producer.Node  `json:"node,omitempty"`
}

func (e EnhancedAlarm) enhance(alarm *producer.Alarm, ctx goka.Context, nodesTopic goka.Table) *EnhancedAlarm {
	enhanced := &EnhancedAlarm{Alarm: alarm}
	log.Printf("enhancing alarm with ID %d and reducion-key %s", alarm.Id, alarm.ReductionKey)
	if c := alarm.NodeCriteria; c != nil {
		node := ctx.Lookup(nodesTopic, producer.GetNodeKey(c))
		if node != nil {
			enhanced.Node = node.(*producer.Node)
		}
	}
	return enhanced
}

// DefineGroup creates a Goke GroupGraph to enhance alarms
func (e EnhancedAlarm) DefineGroup(groupID, alarmsTopic, nodesTopic, targetTopic string) *goka.GroupGraph {
	return goka.DefineGroup(goka.Group(groupID),
		goka.Input(goka.Stream(alarmsTopic), new(producer.AlarmCodec), func(ctx goka.Context, msg interface{}) {
			alarm := msg.(*producer.Alarm)
			enhanced := e.enhance(alarm, ctx, goka.Table(nodesTopic))
			ctx.Emit(goka.Stream(targetTopic), ctx.Key(), enhanced)
		}),
		goka.Output(goka.Stream(targetTopic), new(EnhancedAlarmCodec)),
		goka.Lookup(goka.Table(nodesTopic), new(producer.NodeCodec)),
	)
}

// EnhancedAlarmCodec a Goka Codec to encode/decode an enhanced alarm to/from JSON
type EnhancedAlarmCodec struct{}

// Encode encodes an OpenNMS Alarm object into a JSON byte array
func (c EnhancedAlarmCodec) Encode(value interface{}) (data []byte, err error) {
	data, err = json.Marshal(value)
	return data, err
}

// Decode decodes a JSON byte array into an OpenNMS Alarm object
func (c EnhancedAlarmCodec) Decode(data []byte) (value interface{}, err error) {
	value = &EnhancedAlarm{}
	err = json.Unmarshal(data, value.(*EnhancedAlarm))
	return value, err
}
