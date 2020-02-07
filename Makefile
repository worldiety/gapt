artifactName = gapt

buildDir = build
TOOLSDIR = build/tools
GO = go
GOLANGCI_LINT = ${TOOLSDIR}/golangci-lint
GOLINT = ${TOOLSDIR}/golint
TMP_GOPATH = $(CURDIR)/${TOOLSDIR}/.gopath

all: generate lint test build ## Runs lint, test and build

clean: ## Removes any temporary and output files
	rm -rf ${buildDir}

lint: ## Executes all linters
	${GOLINT}
	${GOLANGCI_LINT} run --enable-all

test: ## Executes the tests
	${GO} test -race ./...

.PHONY: build
build: ## Performs a build and puts everything into the build directory
	${GO} build cmd/main.go -o ${buildDir}/${artifactName}

run: build ## Starts the compiled program
	${buildDir}/${artifactName}

help: ## Shows this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate
generate: ## Executes go generate
	${GO} generate

setup: installGolint installGolangCi ## Installs golint and golangci-lint
	${GO} mod tidy

installGolint:
	GOPATH=${TMP_GOPATH} && go get -u golang.org/x/lint/golint && go install golang.org/x/lint/golint
	cp ${TMP_GOPATH}/bin/golint ${TOOLSDIR}

installGolangCi:
	mkdir -p ${TOOLSDIR}
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TOOLSDIR)/ v1.23.1

.DEFAULT_GOAL := help

