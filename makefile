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
