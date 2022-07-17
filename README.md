# Go bimg Formpost
Multipart formpost microservice based on [`bimg`](https://github.com/h2non/bimg), a very fast image processing library which uses native library [`libvips`](https://libvips.github.io/libvips/) via C bindings.

## Build & Launch
Besides Go 1.16, `libvips` needs to be installed separately. `bimg` apparently works with `libvips` 8.3, but recommends 8.8+.

### Locally
Follow [installation instructions](https://github.com/h2non/bimg#libvips) for `libvips`, e.g.
```bash
$ apt install libvips
```
Run the application at port 8080:
```bash
$ go run .
```

### With Docker
```bash
$ docker build --pull --rm -t go-bimg-formpost:latest .
$ docker run --rm -it -p 8080:8080/tcp go-bimg-formpost:latest
```

### With `docker-compose`
```bash
$ docker-compose up
```

## Usage
Send a `multipart/form-data` request with a single file named `img` to the `/convert` endpoint. Width and height of the target image is set via the URL query parameters `width` and `height` respectively. Optionally, image output `format` can be specified (`jpeg`, `png`, `gif`, `tiff`, `webp`, `heif`, or `avif`)

```bash
$ curl -F img=@/path/to/in.png -o /path/to/out.png http://localhost:8080/convert?width=200&height=200&format=png
```

## See Also
* [Official microservice implementation from `bimg` authors](https://github.com/h2non/imaginary)
* [Available `bimg.Process` options](https://pkg.go.dev/github.com/h2non/bimg#Options)