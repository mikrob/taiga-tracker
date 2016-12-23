BINARY=taiga-tracker

PHONY: all

test:
	go test  -v ./...

get:
	go get

image:
	cp ~/.ssh/id_rsa_jenkins .
	docker build --no-cache -t eu.gcr.io/scalezen/infra/taiga_tracker:0.3 .
	rm -f id_rsa_jenkins

all:
	go build -o ${BINARY} main.go
