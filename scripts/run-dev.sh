#!/bin/bash

trap 'kill -TERM $PID' TERM INT
go run . &
cd client && npm start &

PID=$!
wait $PID
trap - TERM INT
wait $PID
EXIT_STATUS=$?