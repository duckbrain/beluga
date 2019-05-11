FROM gobuffalo/buffalo:v0.14.3 as builder

RUN mkdir -p $GOPATH/src/github.com/duckbrain/beluga/belugad
WORKDIR $GOPATH/src/github.com/duckbrain/beluga/belugad

ADD package.json .
ADD yarn.lock .
RUN yarn install --no-progress
ADD . .
RUN go get $(go list ./... | grep -v /vendor/)
RUN buffalo build --static -o /bin/app

FROM docker/compose:1.24.0

COPY --from=builder /bin/app /bin/belugad
ENV GO_ENV=production
ENV ADDR=0.0.0.0
EXPOSE 3000

CMD /bin/belugad migrate; /bin/belugad
