#GO11MODULES=on
APP?=$(shell basename `pwd`)
REGISTRY?=gcr.io/images
VER=$(shell git describe --abbrev=0 --tag || echo v0.0.0) # `echo xxx`
COMMIT_SHA=$(shell git rev-parse --short HEAD)
LD_FLAGS="-s -w -X main.GitCommit=${COMMIT_SHA} -X main.SemVer=${VER}"

.PHONY: build
## build: build the application
build: clean 
	@echo "Building..."
	@go build  \
	-race \
	-ldflags ${LD_FLAGS} \
	-o ${APP} 
	
dlv-debug: clean
	@echo "Building for delve debug..."
	@go build \
	-ldflags ${LD_FLAGS} \
	-gcflags="all=-N -l" \
	-o ${APP} 

.PHONY: dev
dev: build
	ENV_ABC="value abc" \
	ENV_EDF="value edf" \
	./${APP} 


.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@rm -rf ${APP}

.PHONY: test
## test: runs go test with default values
test:
	JWT_SECRET_KEY="qq@oshlkol" \
	go test -timeout 300s -v -count=1 -race ./...

.PHONY: update
## update: runs go get -u 
update:
	go get -u ./...

.PHONY: custom
## custom: populate the template
custom:
	@rm -rf go.*
	@echo ${APP} > .gitignore
	@echo populated ${APP} > README.md
	@go mod init github.com/datewu/${APP}
	@go build
	@make test
	@make build
## add github secrets settings
	@git add .
	@git commit -am "init custom"
	@git push 

.PHONY: build-tokenizer
## build-tokenizer: build the tokenizer application
build-tokenizer:
	${MAKE} -c tokenizer build

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor
	
# helper rule for deployment
check-environment:
ifndef ENV
	$(error ENV not set, allowed values - `staging` or `production`)
endif

.PHONY: docker-build
## docker-build: builds the stringifier docker image to registry
docker-build: build
	docker build -t ${APP}:${COMMIT_SHA} .

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: check-environment docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'