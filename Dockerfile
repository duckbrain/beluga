FROM golang:1.12 as builder

ENV GO111MODULE=on GOARCH=amd64 GOOS=linux
COPY . /code/
WORKDIR /code/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine

COPY --from=builder /code/beluga /usr/bin/beluga
CMD [ "/usr/bin/beluga", "--help" ]