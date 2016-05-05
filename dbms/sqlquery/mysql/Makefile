REPONAME = monitors
NAME = querymysql
DOCKERIMAGENAME = benchflow/$(REPONAME):$(NAME)
VERSION = dev
GOPATH_SAVE_RESTORE:=`pwd`"/Godeps/_workspace"

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
	docker run -d --name db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=stuff mysql
	docker run -d -ti -e "ENVCONSUL_CONSUL=$(ENVCONSUL_CONSUL)" \
	-e "MYSQL_DB_NAME=stuff" -e "MYSQL_HOST=localhost" \
	-e "MYSQL_PORT=3306" -e "MYSQL_USER=root" -e "MYSQL_USER_PASSWORD=password" \
	-e "BENCHFLOW_CONTAINER_NAME=$(NAME)" -e "BENCHFLOW_MONITOR_NAME=$(NAME)" -e "BENCHFLOW_DATA_NAME=$(NAME)" \
	-p 8080:8080 --name $(NAME) $(DOCKERIMAGENAME)_$(VERSION)

build_and_test_container_local: build_container_local test_container_local
	
rm_container_local:
	docker rm -f -v db
	docker rm -f -v $(NAME)