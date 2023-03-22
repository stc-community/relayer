FROM golang:1.20 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src/app
COPY ./ .

RUN go get -d -v ./...
RUN go install -v ./...
RUN cd basic && make

FROM hub.oneitfarm.com/library/ubuntu AS final
WORKDIR /app
COPY --from=builder /go/src/app/basic/relayer-basic /app/relayer-basic

ENTRYPOINT ["/app/relayer-basic"]
