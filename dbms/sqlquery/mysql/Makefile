REPONAME = monitors
NAME = querymysql
DOCKERIMAGENAME = benchflow/$(REPONAME):$(NAME)
VERSION = dev
GOPATH_SAVE_RESTORE:=$(GOPATH):`pwd`

.PHONY: all build_release 

all: build_release

save_dependencies:
	cd src/cloud/benchflow/$(REPONAME)/ && \
	GOPATH=$(GOPATH_SAVE_RESTORE) godep save ./... && \
	rm -rf ../../../../Godeps/*.* && \
	rm -rf ../../../../Godeps && \
	mv Godeps/ ../../../.. && \
	cd ../../../..

restore_dependencies:
	GOPATH=$(GOPATH_SAVE_RESTORE) godep restore ./...

clean:
	go clean -i ./...
	rm -rf Godeps/_workspace/pkg

build:
	godep go build -o bin/$(NAME) -v ./...

build_release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -ldflags '-s' -o bin/$(NAME) -v ./...

install:
	godep go install -v ./...
	mv bin/$(REPONAME) bin/$(NAME)

test:
	godep go test ./...

build_container_local:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 godep go build -ldflags '-s' -o bin/$(NAME)_linux -v ./...
	docker build -t $(DOCKERIMAGENAME)_$(VERSION) -f Dockerfile.test .
	rm bin/$(NAME)_linux

test_container_local:
	#TODO

rm_container_local:
	#TODO