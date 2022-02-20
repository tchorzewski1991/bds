## Flight Data System

#### Design philosophy & guidelines

- One of the first enginerring decisions we need to make is to decide on how we want to manage third party dependencies.

- It is worth to remember that adding new dependencies is easier then removing them. Keep it simple.

#### Go modules support

- presense of the ```go 1.17``` version inside the ```go.mod``` file has the following meaning:
    - Its primary purpose is to ensure Golang tooling works properly
    - It tells us the project is compatible with everything related to go 1.17 and beyond
    - It tells us the minimum version of compiler we need to use to build the project

- ```go.mod``` file should always be present at the root of the project. Otherwise the Go tooling might get confused and not work properly

- ```GOMODCACHE``` points to the module cache. It is a location on the disk where all of the 3rd party dependencies go. Thanks to that module cache go compiler can build our source code against that source code.
    - ```go clean -modcache``` should remove all cached modules out of your system

- 3rd party dependencies maintenance
    - ```go mod tidy``` - this is the most crucial command for maintenance of our 3rd party packages. This command always walks through our project and ensures all of the source code necessaary to build or project is persisted on the disk in our module cache. We should always run ```go mod tidy``` after extending our codebase with new package.

- Module mirror - mental model behind what happens when we run ```go mod tidy```
    - The part of understaing of what happens when we run ```go mod tidy``` has to do with another env variable called ```GOPROXY```
    - ```GOPROXY``` by default points to ```http://proxy.golang.org``` and its primary responsibility is to direct tooling to the proxy server where all code we want to download lives. The job of the proxy server is to proxy all of the VCSs (Github, Gitlab, ...) and make the module of code we need accessible under one unified location.
    - when we request module without specific tag version ```go mod tidy``` requests list of all of the available version tags. Since we request for the direct dependency ```go mod tidy``` decides to download the latest gratest version of the module by default.

- TODO: Write notes on checksum db
- TODO: Write notes on vendoring

#### Setup for local k8s environment

- Create new kind cluster
    - Add initial config file under ```k8s/kind/kind-config.yaml```
    - In our case ```kind-config.yaml``` will be used to store setup for ports and to make our service available outside of the k8s environment.
    - ```kind-config.yaml``` is used as a part of ```kind create cluster``` command while setting up a new cluster.
- Update ```k8s``` namespace with some base configuration for ```flights-api``` pod which will be common to all of the environments. In the latter steps it will be modified with ```kustomize```.