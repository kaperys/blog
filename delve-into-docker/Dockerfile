FROM golang
WORKDIR /go/src/github.com/kaperys/delve-into-docker-app
EXPOSE 40000 1541

RUN go get github.com/derekparker/delve/cmd/dlv
ADD main.go .

CMD [ "dlv", "debug", "github.com/kaperys/delve-into-docker-app", "--listen=:40000", "--headless=true", "--api-version=2", "--log" ]
