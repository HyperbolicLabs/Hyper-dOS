##
# hyperdos
#
# @file
# @version v1-alpha

build-epitome-image:
	docker build -t epitome ./epitome

run-epitome-image:
	docker run --rm -it epitome

build-and-run-epitome-image: build-epitome-image run-epitome-image

build:
	cd epitome; \
	go build -o ./epitome

run-epitome:
	$(eval export HYPERBOLIC_TOKEN=dummy)
	$(eval export HYPERBOLIC_GATEWAY_URL=dummy)

	cd epitome; \
	go run . -loglevel debug

epitome-help:
	cd epitome; \
	go run . -help

mod-tidy:
	cd epitome; \
	go mod tidy

test:
	cd epitome; \
	go test ./...
