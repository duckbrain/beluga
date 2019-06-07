FROM golang:1.12 as builder

ENV GO111MODULE=on GOARCH=amd64 GOOS=linux
RUN mkdir -p /code
WORKDIR /code/
COPY go.mod go.sum /code/
RUN go mod download
COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app . 

FROM alpine

COPY --from=builder /code/app /usr/bin/beluga
CMD [ "/usr/bin/beluga", "--help" ]