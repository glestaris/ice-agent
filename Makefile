.PHONY: all \
	help \
	deps update-deps \
	docker push-docker \
	test lint \
	clean

ice-agent:
	CGO_ENABLED=0 go build -ldflags "-s -d -w" -o ice-agent .

###### Help ###################################################################

help:
	@echo '    all ................................. builds the grootfs cli'
	@echo '    deps ................................ installs dependencies'
	@echo '    update-deps ......................... updates dependencies'
	@echo '    test ................................ runs tests in Docker'
	@echo '    lint ................................ lint the Go code'
	@echo '    docker .............................. build the Docker image'
	@echo '    push-docker ......................... push the Docker image to Dockerhub'
	@echo '    clean ............................... clean the built artifact'

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

lint:
	./hack/lint

###### Docker #################################################################

docker:
	docker build -t glestaris/ice-agent-ci .

push-docker:
	docker push glestaris/ice-agent-ci

###### Cleanup ################################################################

clean:
	rm -f ice-agent
