FROM golang:1.22 as Build

WORKDIR /go/src/weathercheck

RUN apt update

COPY . .

RUN go get -d -v ./... && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /go/bin/weathercheck ./cmd/weathercheck

FROM alpine:3.18

COPY --from=Build /go/bin/weathercheck /usr/local/bin/weathercheck

ENTRYPOINT ["/usr/local/bin/weathercheck"]