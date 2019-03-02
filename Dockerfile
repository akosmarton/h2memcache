FROM golang:1.12

ENV GO111MODULE=on

WORKDIR /go/src/h2memcache

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install ./...

ENTRYPOINT ["h2memcache"]