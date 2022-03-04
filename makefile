APP=fds

# =============================================================================
# Local development section

run: 
	go run app/services/flights-api/main.go | go run app/services/tools/fmt/main.go

build:
	go build -o $(APP) -ldflags '-X main.build=local' ./app/services/flights-api

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

monitor:
	go install github.com/divan/expvarmon@latest
	expvarmon -ports="4000"

# =============================================================================
# Docker containers section
VERSION := 1.0

all: flights-api

flights-api:
	docker build \
		-f infra/docker/dockerfile.flights-api \
		-t flights-api-amd64:$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg BUILD_REF=$(VERSION) \
		.

docker-down:
	docker rm -f $(docker ps -aq)

docker-clean:
	docker system prune -f

docker-kind-logs:
	docker logs -f $(KIND_CLUSTER)-control-plane

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
	cd infra/k8s/kind/flights-pod; kustomize edit set image flights-api-image=flights-api-amd64:$(VERSION)
	kind load docker-image flights-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/flights-pod/ | kubectl apply -f -

kind-logs:
	kubectl logs -l app=flights --all-containers=true -f --tail=100 | go run app/services/tools/fmt/main.go

kind-restart:
	kubectl rollout restart deployment flights-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=flights