# Go Image Scaling Microservice Example with bimg
Example implementation of a REST microservice based on [`bimg`](https://github.com/h2non/bimg), a very fast image processing library which uses native library [`libvips`](https://libvips.github.io/libvips/) via C bindings.

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
$ docker build --pull --rm -t bimg-rest-example:latest .
$ docker run --rm -it  -p 8080:8080/tcp bimg-rest-example:latest
```

### With `docker-compose`
```bash
$ docker-compose up
```

## Usage
Send a `multipart/form-data` request with a single file named `image` to the `/scale` endpoint. Width and height of the target image is set via the URL query parameters `width` and `height` respectively. Optionally, image output `format` can be specified (`jpeg`, `png`, `gif`, `webp`, `heif`, or `avif`)

```bash
$ curl -F image=@/path/to/in.png -o /path/to/out.png http://localhost:8080/scale?width=200&height=200&format=png
```

Alternatively, there is also `test.html` with HTML `form` and `input`. While experimenting, just edit the `form action` URL.


## See Also
* [Official microservice implementation from `bimg` authors](https://github.com/h2non/imaginary)
* [Available `bimg.Process` options](https://pkg.go.dev/github.com/h2non/bimg#Options)