FROM golang:1.11.1-alpine as build
RUN apk add --no-cache git

# deps
RUN go get -u "github.com/dgraph-io/badger" \
 && go get -u "github.com/gin-gonic/gin" \
 && go get -u "github.com/Depado/ginprom" \
 && go get -u "github.com/jmoiron/sqlx" \
 && go get -u "github.com/lib/pq" \
 && go get -u "github.com/jawher/mow.cli"

RUN go get -u "github.com/graphql-go/graphql" \
    go get -u "github.com/graphql-go/handler"

RUN mkdir -p /go/src/github.com/wcrbrm/badger-example-heartbeat/server
WORKDIR /go/src/github.com/wcrbrm/badger-example-heartbeat/server

# Uploading source files
COPY *.go ./

# Compiling
RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/server

FROM alpine:3.9
COPY --from=build /usr/bin/server /root/
EXPOSE 8092
WORKDIR /root/
VOLUME /badger
CMD ["./server"]
