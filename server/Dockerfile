FROM golang:1.12

WORKDIR /go/src/h2memcache
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV CACHE_SIZE ""
ENV HOST ""
ENV PORT ""
ENV TLS_CERT_FILE ""
ENV TLS_KEY_FILE ""
ENV TLS_CERT_DIR ""
ENV API_KEY ""

CMD ["h2memcache"]