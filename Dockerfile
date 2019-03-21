FROM golang:1.12 as builder

ENV GO111MODULE=on
COPY . /code
WORKDIR /code
RUN go install ./...

FROM alpine

COPY --from=builder /go/bin/beluga* /usr/bin/
