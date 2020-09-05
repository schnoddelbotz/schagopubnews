
BINARY := build/spn

VERSION := $(shell git describe --tags | cut -dv -f2)
DOCKER_IMAGE := schnoddelbotz/schagopubnews
LDFLAGS := -X github.com/schnoddelbotz/schagopubnews/cmd.AppVersion=$(VERSION) -w

ASSETS := handlers/assets.go
GO_SOURCES := */*.go */*/*.go $(ASSETS)
GCP_GO_RUNTIME := go113

GCP_PROJECT := hacker-playground-254920
GCP_REGION := europe-west1

build: $(BINARY)

$(BINARY): $(GO_SOURCES)
	# building schagopubnews
	go build -v -o $(BINARY) -ldflags='-w -s $(LDFLAGS)' ./cli/spn

$(ASSETS):
	test -n "$(shell which esc)" || go get -v -u github.com/mjibson/esc
	go generate handlers/http_handler.go

all_local: clean test build

all_docker: clean test docker_image_prod

release: all_docker docker_image_push


test:
	# golint -set_exit_status ./...
	go fmt ./...
	# golint ./...
	go vet ./...
	go test -ldflags='-w -s $(LDFLAGS)' ./...
	go build cloudfunction.go

coverage: clean
	PROVIDER=MEMORY go test -coverprofile=coverage.out -coverpkg=./... -ldflags='-w -s $(LDFLAGS)' ./...
	go tool cover -html=coverage.out

deploy_gcp: test clean
	gcloud functions deploy Schagopubnews --region=$(GCP_REGION) --runtime=$(GCP_GO_RUNTIME) \
 		--trigger-http --allow-unauthenticated --project=$(GCP_PROJECT) \
 		--set-env-vars=SPN_VERSION=$(VERSION),SPN_REGION=$(GCP_REGION),SPN_PROJECT=$(GCP_PROJECT)

docker_image: clean
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .


docker_image_push:
	docker push $(DOCKER_IMAGE)

docker_run:
	# are you sure docker image is up to date?
	docker run --rm $(DOCKER_IMAGE):latest version

clean:
	rm -f $(BINARY) $(ASSETS) coverage*
