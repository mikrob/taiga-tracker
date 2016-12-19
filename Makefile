BINARY=taiga-tracker

PHONY: all

test:
	go test  -v ./...

get:
	go get

image:
	docker build -t eu.gcr.io/scalezen/infra/taiga_tracker:0.1.0 .

all:
	go build -o ${BINARY} main.go
