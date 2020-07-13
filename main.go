package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agalue/producer-enhancer/model"
	"github.com/burdiyan/kafkautil"
	"github.com/lovoo/goka"
)

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

	var group *goka.GroupGraph
	switch targetKind {
	case "events":
		group = new(model.EnhancedEvent).DefineGroup(groupID, eventsTopic, nodesTopic, targetTopic)
	case "alarms":
		group = new(model.EnhancedAlarm).DefineGroup(groupID, alarmsTopic, nodesTopic, targetTopic)
	default:
		log.Fatalf("invalid target-kind %s. Valid options: 'events' or 'alarms'", targetKind)
	}

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
