#!/bin/bash

echo "Building test container for TaskManager ..."

docker build -t taskmanager-test -f cmd/TaskManager/Dockerfile .