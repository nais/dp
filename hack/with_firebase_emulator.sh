#!/bin/bash

cmd=$1
max_attempts=10
interval=2
emulator_port=6969

gcloud beta emulators firestore start --host-port=localhost:${emulator_port} &> /dev/null &
echo "waiting for firestore emulator to start..."

for attempt in $(seq 1 $max_attempts); do
  echo "attempt ${attempt}/${max_attempts}"
  if netstat -an | grep LISTEN | grep ${emulator_port} &>/dev/null; then
    echo "firestore is running"
    break
  fi
  sleep ${interval}
done

FIRESTORE_EMULATOR_HOST=localhost:${emulator_port} $cmd
pkill -f firestore
