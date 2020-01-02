#!/bin/bash

protoc -I./pb/ ./pb/*.proto --gogoslick_out=plugins=grpc:./pb/
