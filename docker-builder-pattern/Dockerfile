FROM golang
COPY . /go/src/github.com/kaperys/blog/docker-builder-pattern
WORKDIR /go/src/github.com/kaperys/blog/docker-builder-pattern
RUN go get && CGO_ENABLED=0 GOOS=linux go build -o server .

FROM scratch
LABEL maintainer="Mike Kaperys <mike@kaperys.io>"
COPY --from=0 /go/src/github.com/kaperys/blog/docker-builder-pattern/server /opt/kaperys/vision/server
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD html/ /opt/kaperys/vision/html
EXPOSE 8080
WORKDIR /opt/kaperys/vision
ENTRYPOINT [ "./server" ]
