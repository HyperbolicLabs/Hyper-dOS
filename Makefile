build-local-sshbox-image:
	docker build -t local-sshbox -f local/Dockerfile .

minikube-import-images:
	minikube --profile hyperbolic \
	--alsologtostderr image load \
	--remote=false --pull=false 'local-sshbox:latest'

deploy-local-yaml:
	cd local; \
	kubectl apply -f example.yaml


run-sshbox-image:
	docker run -t local-sshbox --rm -d -p '2222:2222' local-sshbox


build-epitome-image:
	docker build -t epitome ./epitome

run-epitome-image:
	docker run --rm -it epitome

build-and-run-epitome-image: build-epitome-image run-epitome-image

build:
	cd epitome; \
	go build -o ./epitome

run-epitome:
	cd epitome; \
	export HYPERBOLIC_GATEWAY_URL='https://api.dev-hyperbolic.xyz' && \
	go run . -loglevel debug -kubeconfig ~/.kube/config

epitome-help:
	cd epitome; \
	go run . -help

mod-tidy:
	cd epitome; \
	go mod tidy

test:
	cd epitome; \
	go test ./...

.PHONY: helm-test
helm-test:
	@cd metadeployment; \
	helm template metadeployment \
		--set ref="dev" \
		.

	@cd gitapps/nvidia-smi; \
	helm template nvidia-smi \
		--set ref="dev" \
		.

	@cd gitapps/epitome; \
	helm template epitome \
		--set ref="dev" \
		.

	@cd gitapps/epitome; \
	helm template epitome \
		--set ref="main" \
		.

	@cd gitapps/pre-pull; \
	helm template pre-pull \
		--set ref="main" \
		.

test-helm-install:
	@cd charts/hyperdos; \
	helm template hyperdos \
		--debug \
		--set ref="dev" \
		.

.PHONY: aider
aider:
	cd hack; \
	bash run-aider.sh
