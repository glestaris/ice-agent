.PHONY: all \
	help \
	deps update-deps \
	test \
	go-vet \
	lint

all:
	GOOS=linux go build -o ice-agent .

###### Help ###################################################################

help:
	@echo '    all ................................. builds the grootfs cli'
	@echo '    deps ................................ installs dependencies'
	@echo '    update-deps ......................... updates dependencies'
	@echo '    test ................................ runs tests in Docker'
	@echo '    go-vet .............................. runs go vet in grootfs source code'
	@echo '    lint ................................ lint the Go code'
	@echo '    docker .............................. build the Docker image'
	@echo '    docker-push ......................... push the built Docker image'

###### Dependencies ###########################################################

deps:
	glide install

update-deps:
	glide update

###### Testing ################################################################

test:
	./hack/run-tests
	./hack/run-tests -i

###### Code quality ###########################################################

go-vet:
	GOOS=linux go vet `go list ./... | grep -v vendor`

lint:
	./hack/lint

###### Docker #################################################################

docker:
	docker build -t glestaris/ice-agent-test .

docker-push:
	docker push glestaris/ice-agent-test
