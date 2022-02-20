APP=fds

# =============================================================================
# Local development section

run: 
	go run main.go

build:
	go build -o $(APP) -ldflags '-X main.build=local'

clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	@echo "  >  Removing $(APP) executable"
	@rm $(APP) 2> /dev/null | true

staticcheck:
	staticcheck -checks=all ./...

tidy:
	go mod tidy
	go mod vendor

# =============================================================================
# Docker containers section
VERSION := 1.0

all: flights-api

flights-api:
	docker build \
		-f infra/docker/dockerfile \
		-t flights-api-amd64:$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg BUILD_REF=$(VERSION) \
		.

# =============================================================================
# Running withing k8s cluster

KIND_CLUSTER := fds-cluster

kind-up:
	kind create cluster \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=flights-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-flights-system:
	kubectl get pods -o wide --watch

kind-load:
	kind load docker-image flights-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/flights-pod/ | kubectl apply -f -

kind-logs:
	kubectl logs -l app=flights --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment flights-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=flights