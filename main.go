package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/agalue/producer-enhancer/model"
	"github.com/agalue/producer-enhancer/protobuf/producer"
	"github.com/burdiyan/kafkautil"
	"github.com/lovoo/goka"
)

func enhanceAlarm(alarm *producer.Alarm, ctx goka.Context, nodesTopic goka.Table) *model.EnhancedAlarm {
	enhanced := &model.EnhancedAlarm{Alarm: alarm}
	log.Printf("enhancing alarm with ID %d and reducion-key %s", alarm.Id, alarm.ReductionKey)
	if c := alarm.NodeCriteria; c != nil {
		var key string
		if c.ForeignSource != "" && c.ForeignId != "" {
			key = fmt.Sprintf("%s:%s", c.ForeignSource, c.ForeignId)
		} else {
			key = fmt.Sprintf("%d", alarm.NodeCriteria.Id)
		}
		node := ctx.Lookup(nodesTopic, key)
		enhanced.Node = node.(*producer.Node)
	}
	return enhanced
}

func main() {
	var bootstrap, groupID, nodesTopic, alarmsTopic, targetTopic string
	flag.StringVar(&bootstrap, "bootstrap", "localhost:9092", "Kafka Bootsrap Server")
	flag.StringVar(&groupID, "group-id", "onms-alarm-group", "Group ID")
	flag.StringVar(&nodesTopic, "nodes-topic", "nodes", "OpenNMS Nodes Topic")
	flag.StringVar(&alarmsTopic, "alarms-topic", "alarms", "OpenNMS Alarms Topic")
	flag.StringVar(&targetTopic, "target-topic", "enhanced-alarms", "The target topic for the enhanced Alarms")
	flag.Parse()

	group := goka.DefineGroup(goka.Group(groupID),
		goka.Input(goka.Stream(alarmsTopic), new(producer.AlarmCodec), func(ctx goka.Context, msg interface{}) {
			alarm := msg.(*producer.Alarm)
			enhanced := enhanceAlarm(alarm, ctx, goka.Table(nodesTopic))
			ctx.Emit(goka.Stream(targetTopic), ctx.Key(), enhanced)
		}),
		goka.Output(goka.Stream(targetTopic), new(model.EnhancedAlarmCodec)),
		goka.Lookup(goka.Table(nodesTopic), new(producer.NodeCodec)),
	)

	processor, err := goka.NewProcessor([]string{bootstrap}, group, goka.WithHasher(kafkautil.MurmurHasher))
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

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
