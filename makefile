APP=bds

# =============================================================================
# Local development section

run: 
	go run app/services/books-api/main.go | go run app/services/tools/fmt/main.go

build:
	go build -o $(APP) -ldflags '-X main.build=local' ./app/services/books-api

clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	@echo "  >  Removing $(APP) executable"
	@rm $(APP) 2> /dev/null | true

lint:
	@golangci-lint run -v -c golangci.yaml

tidy:
	go mod tidy
	go mod verify
	go mod vendor

monitor:
	go install github.com/divan/expvarmon@latest
	expvarmon -ports="4000" -vars="requests,goroutines,errors,panics"

gentoken:
	# -sub=X (user by default) -iss=X (bds-toolset by default) -dur=X (1h by default) -perm=X ("" by default)
	go run app/services/tools/gentoken/main.go

dblab:
	dblab --host localhost --user postgres --pass password --ssl disable --port 5432 --driver postgres

dbmigrate:
	go run app/services/tools/dbmigrate/main.go

# =============================================================================
# Docker containers section
VERSION := 1.0

all: books-api

books-api:
	docker build \
		-f infra/docker/dockerfile.books-api \
		-t books-api-amd64:$(VERSION) \
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

KIND_CLUSTER := bds-cluster

kind-up:
	kind create cluster \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=books-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-books:
	kubectl get pods -o wide --watch --namespace=books-system

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-load:
	cd infra/k8s/kind/books-pod; kustomize edit set image books-api-image=books-api-amd64:$(VERSION)
	kind load docker-image books-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	kustomize build infra/k8s/kind/books-pod/ | kubectl apply -f -

kind-logs:
	#kubectl logs -l app=books --all-containers=true -f --tail=100 | go run app/services/tools/fmt/main.go
	kubectl logs -l app=books --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment books-pod

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	kubectl describe pod -l app=books

kind-access-db:
	kubectl exec -it $(shell kubectl get pods --namespace database-system --output=jsonpath={.items..metadata.name}) --namespace=database-system -- psql -d postgres -U postgres