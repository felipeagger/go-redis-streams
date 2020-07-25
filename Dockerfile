FROM golang:1.13-alpine as builder
ENV GOPATH="$HOME/go"
RUN apk --no-cache add git
WORKDIR $GOPATH/src

COPY . $GOPATH/src

RUN go get -d -v golang.org/x/net/html
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app consumer/*.go


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder $HOME/go/src/app .
CMD ["./app"]  