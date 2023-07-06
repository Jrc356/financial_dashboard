#!/bin/bash

trap 'kill -TERM $PID' TERM INT
go run . &
sleep 5
cd client && npm start &

PID=$!
wait $PID
trap - TERM INT
wait $PID
EXIT_STATUS=$?