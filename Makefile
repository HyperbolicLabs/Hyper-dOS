build-local-sshbox-image:
	docker build -t local-sshbox -f local/Dockerfile .

build-sshbox:
	# DOCKER_BUILDKIT=0 \

	docker build -t sshbox -f images/sshbox/Dockerfile .

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

epitomesh:
	@echo "Starting epitome in sh mode..."
	@cd epitome; go run . -mode sh

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
		--set cascade.hyperdos.ref="dev" \
		.

	@cd gitapps/buffalo/nvidia-smi; \
	helm template nvidia-smi \
		.

	@cd gitapps/epitome; \
	helm template epitome \
		--set cascade.hyperdos.ref="dev" \
		--set cascade.jungleRole.buffalo="true" \
		.

	@cd gitapps/buffalo/pre-pull; \
	helm template pre-pull \
		--set cascade.hyperdos.ref="dev" \
		.

	@cd gitapps/buffalo/instance; \
	helm template instance \
		--set sshPubKeys="pubkey1\npubkey2" \
		--set pubkeyConfig="instance-name" \
		.

	@cd gitapps/buffalo/hyperpool; \
	helm template hyperpool \
		--set enabled="true", \
		--set cascade.hyperpool.models[0].name="test-model" \
		--set cascade.hyperpool.models[0].model="test-model-str" \
		--set cascade.hyperpool.models[0].priority=-1000 \
		--set cascade.hyperpool.models[0].extraArgs[0]="--dtype=half" \
		--set cascade.hyperpool.models[0].resources.requests.cpu=1 \
		--set cascade.hyperpool.models[0].resources.requests.nvidia.com\/gpu=1 \
		--set cascade.hyperpool.models[0].resources.limits.cpu=1 \
		--set cascade.hyperpool.models[0].resources.limits.nvidia.com\/gpu=1 \
		.

test-helm-install:
	@cd charts/hyperdos; \
	helm template hyperdos \
		--debug \
		--set cascade.hyperdos.ref="dev" \
		.

.PHONY: aider
aider:
	cd hack; \
	bash run-aider.sh
