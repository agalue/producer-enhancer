package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agalue/producer-enhancer/model"
	"github.com/agalue/producer-enhancer/protobuf/producer"
	"github.com/burdiyan/kafkautil"
	"github.com/lovoo/goka"
)

func getNodeKey(c *producer.NodeCriteria) string {
	var key string
	if c.ForeignSource != "" && c.ForeignId != "" {
		key = fmt.Sprintf("%s:%s", c.ForeignSource, c.ForeignId)
	} else {
		key = fmt.Sprintf("%d", c.Id)
	}
	return key
}

func enhanceEvent(event *producer.Event, ctx goka.Context, nodesTopic goka.Table) *model.EnhancedEvent {
	enhanced := &model.EnhancedEvent{Event: event}
	log.Printf("enhancing event with ID %d and UEI %s", event.Id, event.Uei)
	if c := event.NodeCriteria; c != nil {
		node := ctx.Lookup(nodesTopic, getNodeKey(c))
		if node != nil {
			enhanced.Node = node.(*producer.Node)
		}
	}
	return enhanced
}

func enhanceAlarm(alarm *producer.Alarm, ctx goka.Context, nodesTopic goka.Table) *model.EnhancedAlarm {
	enhanced := &model.EnhancedAlarm{Alarm: alarm}
	log.Printf("enhancing alarm with ID %d and reducion-key %s", alarm.Id, alarm.ReductionKey)
	if c := alarm.NodeCriteria; c != nil {
		node := ctx.Lookup(nodesTopic, getNodeKey(c))
		if node != nil {
			enhanced.Node = node.(*producer.Node)
		}
	}
	return enhanced
}

func main() {
	var bootstrap, groupID, nodesTopic, eventsTopic, alarmsTopic, targetTopic, targetKind string
	flag.StringVar(&bootstrap, "bootstrap", "localhost:9092", "Kafka Bootsrap Server")
	flag.StringVar(&groupID, "group-id", "onms-enhancer-group", "Group ID")
	flag.StringVar(&nodesTopic, "nodes-topic", "nodes", "OpenNMS Nodes Topic")
	flag.StringVar(&eventsTopic, "events-topic", "events", "OpenNMS Events Topic")
	flag.StringVar(&alarmsTopic, "alarms-topic", "alarms", "OpenNMS Alarms Topic")
	flag.StringVar(&targetTopic, "target-topic", "enhanced", "The target topic for the enhanced objects")
	flag.StringVar(&targetKind, "target-kind", "alarms", "The target kind: 'events' or 'alarms'")
	flag.Parse()

	var input goka.Edge
	switch targetKind {
	case "events":
		input = goka.Input(goka.Stream(eventsTopic), new(producer.EventCodec), func(ctx goka.Context, msg interface{}) {
			event := msg.(*producer.Event)
			enhanced := enhanceEvent(event, ctx, goka.Table(nodesTopic))
			ctx.Emit(goka.Stream(targetTopic), ctx.Key(), enhanced)
		})
	case "alarms":
		input = goka.Input(goka.Stream(alarmsTopic), new(producer.AlarmCodec), func(ctx goka.Context, msg interface{}) {
			alarm := msg.(*producer.Alarm)
			enhanced := enhanceAlarm(alarm, ctx, goka.Table(nodesTopic))
			ctx.Emit(goka.Stream(targetTopic), ctx.Key(), enhanced)
		})
	default:
		log.Fatalf("invalid target-kind %s. Valid options: 'events' or 'alarms'", targetKind)
	}

	group := goka.DefineGroup(goka.Group(groupID),
		input,
		goka.Output(goka.Stream(targetTopic), new(model.EnhancedAlarmCodec)),
		goka.Lookup(goka.Table(nodesTopic), new(producer.NodeCodec)),
	)

	var processor *goka.Processor
	var err error
	var connected = false
	for !connected {
		processor, err = goka.NewProcessor([]string{bootstrap}, group, goka.WithHasher(kafkautil.MurmurHasher))
		if err != nil {
			log.Printf("error creating processor: %v", err)
			log.Println("trying again in 10 seconds...")
			time.Sleep(10 * time.Second)
		} else {
			connected = true
		}
	}
	log.Println("connected.")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		if err = processor.Run(ctx); err != nil {
			log.Printf("error running processor: %v", err)
		}
	}()

	sigs := make(chan os.Signal)
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	}()

	select {
	case <-sigs:
	case <-done:
	}
	cancel()
	<-done
}
