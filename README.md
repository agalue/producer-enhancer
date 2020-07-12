# OpenNMS Kafka Producer Enhancer

A Goka application to generate enhanced version of Alarms with Node data when available.

This solution requires using the OpenNMS Kafka Producer. This feature can export events, alarms, metrics, nodes and edges from the OpenNMS database to Kafka. All the payloads are stored using Google Protobuf.

This repository also contains a Dockerfile to compile and build an image with the tool, which can be fully customized through environment variables.

The `protobuf` directory contains the GPB definitions extracted from OpenNMS source code contains. If any of those files change, make sure to re-generate the protobuf code by using the [build.sh](protobuf/build.sh) command, which expects to have `protoc` installed on your system.

> Future updates will also cover OpenNMS Events.

## Requirements

* `BOOTSTRAP_SERVER` environment variable with Kafka Bootstrap Server (i.e. `kafka01:9092`)
* `NODE_TOPIC` environment variable with the nodes Kafka Topic with GPB Payload (defaults to `nodes`)
* `ALARM_TOPIC` environment variable with the alarms Kafka Topic with GPB Payload (defaults to `alarms`)
* `TARGET_TOPIC` environment variable with the target Kafka Topic to hold the enhanced alarm in JSON (defaults to `enhanced_alarms`)
* `GROUP_ID` environment variable with the Consumer Group ID (defaults to `onms-alarm-group`)

The above environment variables are related with the CLI flags of the application.

## Build

In order to build the application:

```bash
docker build -t agalue/producer-enhancer-go:latest .
docker push agalue/producer-enhancer-go:latest
```

> *NOTE*: Please use your own Docker Hub account or use the image provided on my account.

To build the controller locally for testing:

```bash
export GO111MODULE="on"

go build
./producer-enhancer
```