FROM golang:1.17-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/awdarlan/sinaq

COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v && go build -o /bin/sinaq

FROM alpine

WORKDIR /hrp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/sinaq /awdarlan/sinaq
ENTRYPOINT ["/awdarlan/sinaq"]