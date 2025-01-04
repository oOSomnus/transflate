#!/bin/bash

echo "Building test container for TaskManager ..."

docker build -t taskmanager-test -f cmd/task_manager/Dockerfile .
docker build -t ocrservice-test -f cmd/ocr_service/Dockerfile .
docker build -t translateservice-test -f cmd/translate_service/Dockerfile .