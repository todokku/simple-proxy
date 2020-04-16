FROM golang:1.14 as builder

COPY main.go .

ARG go_arch
RUN GOARCH=$go_arch go build -o proxy main.go

# arm v7 3.11
FROM alpine@sha256:c5ea49127cd44d0f50eafda229a056bb83b6e691883c56fd863d42675fae3909

COPY --from=builder /go/proxy /proxy

ENTRYPOINT ["/proxy"]
