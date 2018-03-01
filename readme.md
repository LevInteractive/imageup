# ImageUp for Google Cloud Platform (Storage)

It's recommended this be run as a private microservice (most likely within a
Kubernetes cluster) because it does not handle any type of authentication. That
should be done by the application interfacing with this service.

### Usage

You can use docker or just create a binary and run it.

```
docker run -it --rm levinteractive/imageup \
  -e GOOGLE_APPLICATION_CREDENTIALS=/path/to/servicefile.json \
  -e BUCKET_ID="my-bucket" \
  -e SERVER_PORT="8080"
```

### Environmental Variables

* `BUCKET_ID` default: null (required)
* `CACHE_MAX_AGE` default: 86400
* `SERVER_PORT` default: 31111
* `CORS` default: "*"

### API

##### Remove file(s) from storage

Send a `DELETE` to `/` with the following argument.

* `files` - A comma-delimited list of files. "file1, file2, etc"

*Example Request:*

```shell
curl -X DELETE \
  http://localhost:31111/ \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'files=foobar-82f100f-4bfd-4671-a0c1-9fc662782.jpg,foobar-38ef100f-4bfd-4671-a0c1-9fc667f2d924.jpg'
```


##### Add files to storage

Send a `POST` to `/` with the following arguments.

* `file` - Binary file data)
* `sizes` - Array of file configurations to return

*sizes:*

The `file` will be converted to all of the dimensions specified in the `sizes`
argument. If the array below was sent, two images would be saved and returned.

```json
[
  {
    "name": "large",
    "width": 1000,
    "height": 1000,
    "fill": false
  },
  {
    "name": "thumb",
    "width": 500,
    "height": 500,
    "fill": true
  }
]
```

*Example Request:*

```shell
curl -X POST \
  http://localhost:31111/ \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F file=@_DSC1814.jpg \
  -F 'sizes=[{"name": "foobar", "width":100, "height": 100, "fill": false}]'
```

### Response

The response objects will always be in the same order they were sent.

```json
[
  {
    "name": "large",
    "url": "https://the-public-url",
    "width": 500,
    "height": 500,
    "fill": true
  },
  {
    "name": "thumb",
    "url": "https://the-public-url",
    "width": 500,
    "height": 500,
    "fill": true
  }
]
```

### Errors

The `code` will be a relevant [http code](https://golang.org/pkg/net/http/#pkg-constants).

```json
{
    "code": 405,
    "message": "invalid"
}
```
