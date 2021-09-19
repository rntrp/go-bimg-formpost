FROM golang:1.16-alpine
# Add common fonts for proper SVG text rendering:
RUN apk add --no-cache fontconfig ghostscript-fonts ttf-liberation \
        ttf-dejavu font-noto font-noto-emoji ttf-font-awesome \
        msttcorefonts-installer && update-ms-fonts && fc-cache -f
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY internal ./internal
# Generate native Go bindings for vips, remove build cruft afterwards:
RUN apk add --no-cache --virtual .build vips-dev build-base \
    && go mod download \
    && go test ./... \
    && go build -o /bimg-rest \
    && apk del .build \
    && apk add --no-cache vips \
    && rm -rf /tmp/*
EXPOSE 8080
CMD ["/bimg-rest"]
