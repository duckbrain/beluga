FROM golang:1.14

ENV GO111MODULE=on
WORKDIR /code/
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go env
RUN go build -ldflags "-linkmode external -extldflags -static" -a -o /app ./cmd/beluga

FROM alpine

RUN apk add --no-cache bash docker docker-compose
COPY --from=0 /app /usr/bin/beluga
CMD [ "/usr/bin/beluga", "--help" ]
