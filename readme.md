# ImageUp for Google Cloud Platform (Storage)

It's recommended this be run as a private microservice (most likely within a
Kubernetes cluster) because it does not handle any type of authentication. That
should be done by the application interfacing with this service.

### Environmental Variables

* `BUCKET_ID` default: null (required)
* `CACHE_MAX_AGE` default: 86400
* `SERVER_PORT` default: 31111
* `CORS` default: "*"

### `DELETE` Request


### `POST` Body Arguments

* `file` - Binary file data)
* `sizes` - Array of file configurations to return

##### sizes value example

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

### Response

The response objects will always be in the same order they were sent.

```json
[
  {
    "name": "large",
    "url": "[ public google bucket url ]",
    "width": 500,
    "height": 500,
    "fill": true
  },
  {
    "name": "thumb",
    "url": "[ public google bucket url ]",
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
