#!/bin/bash

type protoc >/dev/null 2>&1 || { echo >&2 "protoc required but it's not installed; aborting."; exit 1; }

protoc --proto_path=./ --go_out=./ opennms-kafka-producer.proto
