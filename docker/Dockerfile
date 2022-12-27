FROM golang:1.19-alpine AS builder

COPY . /build
WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /petly-app cmd/main.go

FROM alpine:latest
COPY --from=builder /petly-app /petly-app
EXPOSE 8080

ENTRYPOINT ["/petly-app", "serve"]
