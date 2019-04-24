FROM golang:alpine AS build-env

RUN apk add --no-cahe git
WORKDIR /go/src/github.com/nomeaning777/simple-webdav-server
COPY server.go . 
RUN go get
RUN go build && cp simple-webdav-server /bin

FROM alpine:3.9
COPY --from=build-env /bin/simple-webdav-server /usr/local/bin/simple-webdav-server
WORKDIR /storage
CMD /usr/local/bin/simple-webdav-server