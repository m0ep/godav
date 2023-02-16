# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

RUN mkdir data

COPY data/test.txt ./data/
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /godav

EXPOSE 80

CMD [ "/godav" ]
