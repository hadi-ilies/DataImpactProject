FROM golang:latest

WORKDIR /go/src/DataImpactProject/server

ADD . /go/src/DataImpactProject

RUN go mod vendor

CMD go run cmd/*.go

EXPOSE 8081
