FROM golang:1.9 AS builder

RUN apt-get update && apt-get install -y upx

ARG PACKAGE

RUN mkdir -p /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

ARG VERSION

ENV LD_FLAGS="-w -X main.Version=${VERSION}"
ENV CGO_ENABLED=0

COPY . /go/src/${PACKAGE}
RUN go get -d -v .
RUN go install -a -v -tags netgo -ldflags "${LD_FLAGS}" .
RUN upx -9 /go/bin/panik

FROM scratch
MAINTAINER Robert Jacob <xperimental@solidproject.de>
EXPOSE 8080

COPY --from=builder /go/bin/panik /panik
ENTRYPOINT ["/panik"]