FROM golang:1.16-alpine3.14 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY internal ./internal
RUN apk add --no-cache build-base vips-dev \
    && go mod download \
    && go test ./... \
    && go build -o /bimg-rest

FROM alpine:3.14
# Add common fonts for proper SVG text rendering:
RUN apk add --no-cache fontconfig ghostscript-fonts ttf-liberation \
        ttf-dejavu font-noto font-noto-emoji ttf-font-awesome \
        msttcorefonts-installer && update-ms-fonts && fc-cache -f \
    && apk add --no-cache vips
COPY --from=builder /bimg-rest ./
EXPOSE 8080
CMD ["/bimg-rest"]
