package model

import (
	"context"
	"testing"

	"github.com/agalue/producer-enhancer/protobuf/producer"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"gotest.tools/v3/assert"
)

func TestEventProcessor(t *testing.T) {
	// Build and run processor
	tt := tester.New(t)
	group := new(EnhancedEvent).DefineGroup("test", "events", "nodes", "enhanced")
	proc, err := goka.NewProcessor(nil, group, goka.WithTester(tt))
	assert.NilError(t, err)
	go proc.Run(context.Background())
	// Create tracker for the enhanced topic
	mt := tt.NewQueueTracker("enhanced")
	// Create node
	node := &producer.Node{
		Id:       1,
		Label:    "test01",
		Location: "Default",
	}
	tt.SetTableValue(goka.Table("nodes"), "1", node)
	// Create event
	event := &producer.Event{
		Id: 10,
		NodeCriteria: &producer.NodeCriteria{
			Id: 1,
		},
		Uei:         "uei.opennms.org/test",
		Description: "This is a test",
		LogMessage:  "This is a test",
	}
	tt.Consume("events", "", event)
	// Verify enhanced event
	if _, value, ok := mt.Next(); ok {
		enhanced := value.(*EnhancedEvent)
		assert.Assert(t, enhanced.Event != nil)
		assert.Equal(t, "This is a test", enhanced.Event.LogMessage)
		assert.Assert(t, enhanced.Node != nil)
		assert.Equal(t, "test01", enhanced.Node.Label)
	} else {
		t.Fail()
	}
}
