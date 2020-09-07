
BINARY := build/spn

VERSION := $(shell git describe --tags | cut -dv -f2)
DOCKER_IMAGE := schnoddelbotz/schagopubnews
LDFLAGS := -X github.com/schnoddelbotz/schagopubnews/cmd.AppVersion=$(VERSION) -w

ASSETS := handlers/assets.go
GO_SOURCES := */*.go */*/*.go $(ASSETS)
UI_SOURCE_ROOT := ../schagopubnews-ui
UI_SOURCE_DIST := $(UI_SOURCE_ROOT)/dist
UI_SOURCES := $(UI_SOURCE_DIST)/index.html $(wildcard $(UI_SOURCE_DIST)/assets/*)
GCP_GO_RUNTIME := go113
GCP_PROJECT := hacker-playground-254920
GCP_REGION := europe-west1

# TODO:
# SPN_UI_URI=...
# SPN_API_URI=...
# pass to ember via on the fly adjusted index.html

build: $(BINARY)

$(BINARY): $(GO_SOURCES)
	# building schagopubnews
	go build -v -o $(BINARY) -ldflags='-w -s $(LDFLAGS)' ./cli/spn

$(UI_SOURCE_ROOT):
	cd .. && git clone git@github.com:schnoddelbotz/schagopubnews-ui.git

$(UI_SOURCE_DIST)/index.html: $(UI_SOURCE_ROOT)
	# echo trying to build schagopubnews-ui in $(UI_SOURCE_ROOT)
	make -C $(UI_SOURCE_ROOT) clean dist

$(ASSETS): $(UI_SOURCES)
	test -n "$(shell which esc)" || go get -v -u github.com/mjibson/esc
	esc -prefix $(UI_SOURCE_ROOT) -pkg handlers -o $(ASSETS) -private $(UI_SOURCE_ROOT)

all_local: clean test build

all_docker: clean test docker_image_prod

release: all_docker docker_image_push

test: $(GO_SOURCES)
	# golint -set_exit_status ./...
	go fmt ./...
	# golint ./...
	go vet ./...
	# running tests...
	go test -ldflags='-w -s $(LDFLAGS)' ./...
	# trying to build cloud function for test purposes
	go build cloudfunction.go

serve: $(BINARY)
	./build/spn serve

coverage: clean
	PROVIDER=MEMORY go test -coverprofile=coverage.out -coverpkg=./... -ldflags='-w -s $(LDFLAGS)' ./...
	go tool cover -html=coverage.out

deploy_gcp: clean test
	# TODO: Note CFN Name and executed function CAN differ!
	gcloud functions deploy SPN --region=$(GCP_REGION) --runtime=$(GCP_GO_RUNTIME) \
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
	make -C $(UI_SOURCE_ROOT) clean

realclean: clean
	-make -C $(UI_SOURCE_ROOT) realclean
	rm -rf schagopubnews-ui