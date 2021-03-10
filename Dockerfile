FROM golang:1.16.0-alpine3.13 as builder

WORKDIR /workdir
COPY . /workdir

RUN mkdir -p target && go build -o target/hare main.go

FROM alpine:3.13
COPY --from=builder /workdir/target/hare /usr/bin/hare
RUN chmod go+rwx /usr/bin/hare
ENTRYPOINT ["hare"]
