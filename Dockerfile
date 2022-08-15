FROM golang:1.18-alpine3.16 AS builder
WORKDIR /app
COPY go.mod go.sum main.go ./
COPY internal ./internal
RUN apk add --no-cache build-base vips-dev upx \
    && go mod download \
    && go test ./... \
    && go build -ldflags="-s -w" -o /go-bimg-formpost \
    && upx --best --lzma /go-bimg-formpost

FROM alpine:3.16
RUN apk add --no-cache vips-poppler ttf-liberation
COPY --from=builder /go-bimg-formpost ./
EXPOSE 8080
ENTRYPOINT ["/go-bimg-formpost"]
