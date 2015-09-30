FROM golang
 
RUN apt-get update \
	&& apt-get -y install mysql-client-5.5 \
    && apt-get -y install p7zip \
    && apt-get clean

ADD src/cloud/benchflow/monitors /go/src/cloud/benchflow/monitors

RUN go get github.com/go-sql-driver/mysql
RUN go install cloud/benchflow/monitors

ENV MYSQL_HOST db
ENV MYSQL_USER root
ENV MYSQL_USER_PASSWORD PASSWORD
ENV MYSQL_DB_NAME stuff

ENTRYPOINT /go/bin/monitors
 
EXPOSE 8080
