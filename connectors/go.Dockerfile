FROM golang:1.18

COPY ./vance /tmp/vance
COPY ./cdk-go /tmp/cdk-go

RUN go build -o /tmp/vance/bin /tmp/vance/connectors/database/mongodb-sink/cmd/main.go