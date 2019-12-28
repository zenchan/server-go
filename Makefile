.PHONY: gateway

ALL_SERVER=gateway

all: $(ALL_SERVER)

gateway:
	go install ./gateway
