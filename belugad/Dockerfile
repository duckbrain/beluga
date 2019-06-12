FROM gobuffalo/buffalo:v0.14.3 as builder

RUN mkdir -p /code
WORKDIR /code/

ENV GO111MODULE=on GOARCH=amd64 GOOS=linux
ADD package.json yarn.lock go.mod go.sum /code/
RUN yarn install --no-progress && go mod download
ADD . /code/
RUN buffalo build --static -o /bin/app

FROM docker/compose:1.24.0

ENV GO_ENV=production
ENV ADDR=0.0.0.0
EXPOSE 3000
COPY --from=builder /bin/app /bin/belugad

ENTRYPOINT []
CMD /bin/sh -c '/bin/belugad migrate; /bin/belugad'
