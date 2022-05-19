BINARY_NAME=esp32-mqtt
FLAGS=-ldflags "-s -w"

all: build

build:
	env CGO_ENABLED=0 go build ${FLAGS} -o ${BINARY_NAME} main.go

run:
	env CGO_ENABLED=0 go build ${FLAGS} -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

.PHONY: all build