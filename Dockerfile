FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY src/*.go ./
# Add common fonts for proper SVG text rendering:
RUN apk add --no-cache fontconfig ghostscript-fonts ttf-liberation \
        font-noto font-noto-emoji \
        msttcorefonts-installer && update-ms-fonts && fc-cache -f
# Generate native Go bindings for vips, remove build cruft afterwards:
RUN apk add --no-cache --virtual .build vips-dev build-base \
    && go mod download \
    && go test \
    && go build -o /bimg-rest \
    && apk del .build \
    && apk add --no-cache vips \
    && rm -rf /tmp/*
EXPOSE 8080
CMD ["/bimg-rest"]
