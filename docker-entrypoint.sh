#!/bin/sh -e

exec /producer-enhancer \
  -bootstrap "${BOOTSTRAP_SERVERS}" \
  -nodes-topic "${NODES_TOPIC-nodes}" \
  -alarms-topic "${ALARMS_TOPIC-alarms}" \
  -target-topic "${TARGET_TOPIC-enhanced-alarms}" \
  -group-id "${GROUP_ID-onms-alarm-group}"
