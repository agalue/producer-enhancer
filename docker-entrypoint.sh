#!/bin/sh -e

exec /producer-enhancer \
  -bootstrap "${BOOTSTRAP_SERVERS}" \
  -nodes-topic "${NODES_TOPIC-nodes}" \
  -events-topic "${EVENTS_TOPIC-events}" \
  -alarms-topic "${ALARMS_TOPIC-alarms}" \
  -target-topic "${TARGET_TOPIC-enhanced}" \
  -target-kind "${TARGET_KIND-alarms}" \
  -group-id "${GROUP_ID-onms-enhancer-group}"
