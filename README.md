# Image Conversion Formpost Microservice
Multipart formpost microservice based on [`bimg`](https://github.com/h2non/bimg), a very fast image processing library which uses native library [`libvips`](https://libvips.github.io/libvips/) via C bindings.

## Build & Launch
Besides Go 1.20, `libvips` needs to be installed separately. `bimg` apparently works with `libvips` 8.3, but recommends 8.8+.

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

### Endpoint Overview
| path | method | description |
|---|---|---|
| `/` | `GET` | Returns a simple HTML page with form uploads for testing the `/convert` endpoint. |
| `/index.html` | `GET` | Same as `/`. |
| `/live` | `GET` | Liveness probe endpoint. Returns HTTP 200 while the application is running. |
| `/metrics` | `GET` | Returns application metrics in Prometheus text-based format collected by the [Prometheus Go client library](https://github.com/prometheus/client_golang). The endpoint is disabled by default; enabled with `BIMG_FORMPOST_ENABLE_PROMETHEUS=true` (see [Environment Variables](#environment-variables) below). |
| `/convert` | `POST` | Pivotal endpoint within the whole application. Accepts a single image per `multipart/form-data` request. Converts the image to the requested image format. Target dimensions are specified as part of the URL query string. See [Parameters](#parameters) below for further details. |
| `/shutdown` | `POST` | Initiates graceful shutdown of the application when triggered by a POST request with an arbitrary payload. Returns HTTP 204 after the shutdown process has been started. The endpoint is disabled by default; enabled with `BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT=true` (see [Environment Variables](#environment-variables) below). |

### Environment Variables
The application supports configuration via environment variables or a `.env` file. Environment variables have higher priority.

| variable | default | description |
|---|---|---|
| `BIMG_FORMPOST_ENV` | `development` | Currently, this setting only affects the [`.env` file precedence](https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use), no actual distinction between execution environments is made. Possible values are `development`, `test` and `production`. |
| `BIMG_FORMPOST_ENV_DIR` | _empty_ | Path to directory containing the `.env` file. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, `.env` file is read from the application folder. |
| `BIMG_FORMPOST_TCP_ADDRESS` | `:8080` | Application TCP address as described by the Golang's `http.Server.Addr` field, most prominently in form of `host:port`. See also [`net` package docs](https://pkg.go.dev/net). |
| `BIMG_FORMPOST_TEMP_DIR` | _OS temp folder_ | Path to directory, where applications's temporary _output_ files are managed. Absolute paths or paths relative to the application folder are possible. If the variable is left _empty_, operating system's default temporary directory is used. Please note that the application always writes _input_ multipart data to the OS' temporary directory, regardless of the `BIMG_FORMPOST_TEMP_DIR` value, unless `BIMG_FORMPOST_MEMORY_BUFFER_SIZE` is negative. This limitation is inflicted by the `mime/multipart` Go API and cannot be feasibly altered (yet). |
| `BIMG_FORMPOST_MAX_REQUEST_SIZE` | `-1` | Maximum size of a multipart request in bytes, which can be processed by the application. Note that the request size amounts to the entire HTTP request body, including multipart boundary delimiters, content disposition headers, line breaks and, indeed, the actual payload. Decent clients which send the `Content-Length` header also enjoy fail-fast behavior if that value exceeds the provided maximum. Either way, the application counts bytes during upload and returns `413 Request Entity Too Large` as soon as the limit is exceeded. By default no request size limit is set. |
| `BIMG_FORMPOST_MEMORY_BUFFER_SIZE` | `10485760` | Number of bytes stored in memory when uploading multipart data. If the payload size is exceeding this number, then the remaining bytes are dumped onto the filesystem. Accordingly, the size of `0` prompts the application to always write all bytes to a temporary file, whereas any negative value such as `-1` will prevent the application from hitting the filesystem and retain the whole request payload in memory. Keep in mind that Go `mime/multipart` always adds 10 MiB on top of this value for the "non-file parts", i.e. boundaries etc., hence the actual minimum is 10 MiB plus 1 byte. The default value is therefore effectively 20 MiB. |
| `BIMG_FORMPOST_ENABLE_PROMETHEUS` | `false` | Expose application metrics via [Prometheus](https://prometheus.io) endpoint `/metrics`. |
| `BIMG_FORMPOST_ENABLE_SHUTDOWN_ENDPOINT` | `false` | Enable shutdown endpoint under `/shutdown`. A single POST request with arbitrary payload to this endpoint will cause the application to shutdown gracefully. |
| `BIMG_FORMPOST_SHUTDOWN_TIMEOUT` | `0s` | Specifies amount of time to wait before ongoing requests are forcefully cancelled in order to perform a graceful shutdown. A zero value lets the application wait indefinetely for all requests to complete. At least one time unit must be specified, e.g. `45s` or `5m15s123ms`. See [Go `time.ParseDuration` format](https://pkg.go.dev/time#ParseDuration) for further details. |

### Parameters
| parameter | mandatory | value |
|---|:---:|---|
| `width` | yes | Min: `1`; Max: `65500` (theoretical). |
| `height` | yes | Min: `1`; Max: `65500` (theoretical). |
| `format` | yes | `jpg`, `jpeg`, `png`, `gif`, `tif`, `tiff`, `webp`, `heif`, `heic`, `avif`. |
| `quality` | no | From `1` to `100` for JPEG (default `95`); `0` to `9` for PNG (default `6`). |
| `resize` | no | Image resizing mode for cases when the source document has a different aspect ratio than specified by the target dimensions. |
| `resample` | no | Resampling algorithm, affects both quality and performance. |
| `background` | no | Background color for flattening images with alpha background, specified by an RGB hex string (default `000000`). |

#### Image Resizing Modes
| value | description |
|---|---|
| `fit` | Image is scaled within the specified width-height box ("fit") preserving the aspect ratio of the image. Scaled images with an aspect ratio differing from the specified dimensions have therefore either less width or height. **Note:** if the source image is smaller than the targeted bounding box, the resulting image will retain its source dimensions. Use `fit-upscale` in case the source image should be upscaled. **`fit` is the default mode.** |
| `fit-upscale` | Same as `fit`, but upscales the source image so either image width or height fits the targeted bounding box. |
| `fit-upscale-black` | Same as `fit`, but the left out pixels are filled with black color. Such images have horizontal or vertical black bars on the opposite sides of the image, a practice also known as "letterboxing" or "pillarboxing", respectively. |
| `fit-upscale-white` | Same as `fit-upscale-black`, but the bars are now white. |
| `fill` | Image is scaled to fill the entire width-height box preserving the aspect ratio of the image. If the image has a different aspect ratio, than specified by the target dimensions, then the image is cropped around the center of the image from both sides. Taller images are equally cropped at the top and bottom, accordingly wider ones are cropped left and right. |
| `fill-north` | Same as `fill`, but taller images are cropped from the top, i.e. northern side, hence the outsized bottom part is left out. Wider images are cropped from the center. This option may come in handy when processing both portrait and landscape oriented documents. |
| `fill-east` | Same as `fill`, but taller images are cropped from the right, i.e. eastern side, hence the left part is left out. Taller images are cropped vertically at the center. |
| `fill-south` | Same as `fill-north`, but taller images are cropped from the bottom, i.e. the southern side. |
| `fill-west` | Same as `fill-east`, but wider images are cropped from the left. |
| `stretch` | The image is stretched to fill the target width and height without preserving the aspect ratio. |

#### Resampling Algorithms
[Wikipedia](https://en.wikipedia.org/wiki/Image_scaling) and [ImageMagick 6 Docs](https://legacy.imagemagick.org/Usage/filter/) are a good starting point for understanding the differences between particular algorithms and performance implications.

| value | description |
|---|---|
| `bicubic` | Bicubic interpolation (**default**) |
| `bilinear` | Bilinear interpolation |
| `nohalo` | NoHalo interpolation, as used in GEGL/GIMP ([further](https://legacy.imagemagick.org/Usage/filter/nicolas/) [reading](https://dl.acm.org/doi/pdf/10.1145/1557626.1557657)) |
| `nearest` | Nearest-neighbor interpolation |

## See Also
* [Official microservice implementation from `bimg` authors](https://github.com/h2non/imaginary)
* [Available `bimg.Process` options](https://pkg.go.dev/github.com/h2non/bimg#Options)