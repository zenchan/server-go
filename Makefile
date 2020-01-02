.PHONY: proto gateway lobby

ALL_SERVER=gateway lobby

all: proto $(ALL_SERVER)

proto:
	cd ./proto && ./gen.sh

gateway:
	go install ./gateway

lobby:
	go install ./lobby
