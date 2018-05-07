build: build-python
	docker pull golang:1
build-go:
	go build corti.ai/queue

build-python:
	cd python; docker build -t corticlient .
