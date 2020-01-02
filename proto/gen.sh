#!/bin/bash

if [ ! -d "pb" ]; then
    mkdir pb
fi
rm -f ./pb/*

protoc -I./proto/ ./proto/*.proto --gogoslick_out=./pb/
