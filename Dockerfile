FROM golang:1.15
MAINTAINER mlabouardy <mohamed@labouardy.com>
WORKDIR /go/src/github.com/mlabouardy/swaggymnia/
COPY . .

RUN go mod init
RUN GO111MODULE=on go get -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o swaggymnia .
ENTRYPOINT ["./swaggymnia"]
