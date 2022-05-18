## Books Data System

#### Design philosophy & guidelines

- One of the first engineering decisions we need to make is to decide on how we want to manage third party dependencies.

- It is worth to remember that adding new dependencies is easier than removing them. Keep it simple.

#### Go modules support

- presence of the ```go 1.17``` version inside the ```go.mod``` file has the following meaning:
    - Its primary purpose is to ensure Golang tooling works properly
    - It tells us the project is compatible with everything related to go 1.17 and beyond
    - It tells us the minimum version of compiler we need to use to build the project

- ```go.mod``` file should always be present at the root of the project. Otherwise, the Go tooling might get confused and not work properly

- ```GOMODCACHE``` points to the module cache. It is a location on the disk where all the 3rd party dependencies go. Thanks to that module cache go compiler can build our source code against that source code.
    - ```go clean -modcache``` should remove all cached modules out of your system

- 3rd party dependencies maintenance
    - ```go mod tidy``` - this is the most crucial command for maintenance of our 3rd party packages. This command always walks through our project and ensures all the source code necessary to build or project is persisted on the disk in our module cache. We should always run ```go mod tidy``` after extending our codebase with new package.

- Module mirror - mental model behind what happens when we run ```go mod tidy```
    - The part of understanding of what happens when we run ```go mod tidy``` has to do with another env variable called ```GOPROXY```
    - ```GOPROXY``` by default points to ```http://proxy.golang.org``` and its primary responsibility is to direct tooling to the proxy server where all code we want to download lives. The job of the proxy server is to proxy all the VCSs (Github, Gitlab, ...) and make the module of code we need accessible under one unified location.
    - when we request module without specific tag version ```go mod tidy``` requests list of all of the available version tags. Since we request for the direct dependency ```go mod tidy``` decides to download the latest gratest version of the module by default.

- TODO: Write notes on checksum db
- TODO: Write notes on endorsing

#### Setup for local k8s environment

- Create new kind cluster
    - Add initial config file under ```k8s/kind/kind-config.yaml```
    - In our case ```kind-config.yaml``` will be used to store setup for ports and to make our service available outside of the k8s environment.
    - ```kind-config.yaml``` is used as a part of ```kind create cluster``` command while setting up a new cluster.
- Update ```k8s``` namespace with some base configuration for ```books-api``` pod which will be common to all of the environments. In the latter steps it will be modified with ```kustomize```.

#### Project layers, policies and guidelines

- Programming mode - is about getting code to work
- Engineering mode - is about getting code to the place where it can be maintained, managed and debugged
- Project layout
    - There is no consensus in Go community about the right project structure. From one point of view this approach makes Go projects extremely flexible to fulfill concrete business requirements. From another point of view this approach makes a lot of inconsistencies and opinions regarding 'the right' way to do things.
    - It takes time and refactoring steps to find the right way to structure your project.
    - One interesting approach might be **layering**:
        - App layer - App layer is on the top of our food chain. The primary responsibility of this layer is to start-up and shutdown the service. It's about getting external input and providing external output. When it comes to REST based APIs the external input/output comes in the form of HTTP **handlers**. When it comes to tooling, external input/output comes in the form of **terminal** stdin/stdout. Code in the app layer is extremely specific to the business case we are trying to fullfil. We shouldn't import anything between packages within this layer. It represents different apps we are building. There are 2 types of applications that we can build: services and tooling.
            - Service represents different APIs we are building, e.x:
                - ```services/books-api``` - business app
                - ```services/users-api``` - business app
                - ```services/metrics``` - sidecar app
            - Tooling represents apps that might not be strictly connected to our business, but might be helpful with maintenance, e.x admin tools:
                - ```tooling/admin/commands/migrate```
                - ```tooling/admin/commands/adduser```
                - ```tooling/admin/commands/gentoken```
        - Business layer - business layer represents the core components necessary to solve business problem in front of us. It might be the composition of more specific layers aimed to solve concrete business problem. It's good to have the following structure:
            - Core - it should be treated as a business layer API. Let's assume we have a functionality responsible for generating books stats. In order to generate a books stats we need to bring everything from data layer. The core layer will leverage the data layer to do the higher level business logic. Anything that expects multiple data layer calls should be leveraged inside the core layer.
            - Data - it should be treated as a raw database layer API. It may contain CRUD Api for your database access objects (DAOs). It is not a problem when app layer is bypassing core layer in order to directly access data layer CRUD operations.
            - System - it should contain packages that will help us resolve business problem we have, but they are system oriented (not specific). Ex. auth, database, metrics, validate
            - Web - it should contain code specific to the web applications we are trying to build, e.x middleware for cors, auth, docs, logger. While trying to organize code under web layer it's worth to consider API versioning quite early. Example package organization:
                - v1
	                - middleware
		                - auth
		                - cors
		                - logger
		                - errors
		                - metrics
                    - v1
                - Why the middleware presented above is not in foundation layer? We might want to change logging, errors or auth  middleware between different API versions.  In this particular case the following layering saves us from making API specific changes inside foundational packages.

        - Foundation layer - this layer contains all the foundational code. Foundational code is all the code that is not specifically tied to any business logic. Foundational packages should be treated a bit like standard libraries. They shouldn't log or wrap any errors. External dependencies makes them less reusable. Any code from app/business layers can communicate directly with packages from that layer. 
        - Vendor layer - this layer contains all the 3rd party code that our app will be using
        - Infra layer - this layer contains all the infrastructure related code and files for k8s | terraform | gcl | docker | etc 
