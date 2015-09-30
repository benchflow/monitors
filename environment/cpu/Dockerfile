FROM golang
 
ADD src/cloud/benchflow/monitors /go/src/cloud/benchflow/monitors

RUN go get github.com/fsouza/go-dockerclient
RUN go install cloud/benchflow/monitors

ENTRYPOINT /go/bin/monitors
 
EXPOSE 8080
