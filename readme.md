# ImageUp for Google Cloud Platform (Storage)

It's recommended this be run as a private microservice (most likely within a
Kubernetes cluster) because it does not handle any type of authentication. That
should be done by the application interfacing with this service.

* [Usage](#usage)
* [API](docs/api.md)
* [Kubernetes example](examples/k8s)
* [Express app example](examples/node)

## Usage

The easiest way to use this is to simply pull and run it using docker. Note that
`GOOGLE_APPLICATION_CREDENTIALS` is only required if you aren't running this in
a Google Cloud environment. Otherwise, it's already set by default.

```
docker run -it --rm levinteractive/imageup \
  -e GOOGLE_APPLICATION_CREDENTIALS=/path/to/servicefile.json \
  -e BUCKET_ID="my-bucket" \
  -e SERVER_PORT="8080"
```

Alternatively, you can download this repo and build the binary for your
respective arch. I didn't include binaries, but could if there is any demand.

## Environmental Variables

* `BUCKET_ID` default: null (required)
* `CACHE_MAX_AGE` default: 86400
* `SERVER_HOST` default: localhost
* `SERVER_PORT` default: 31111
* `CORS` default: "*" (wildcard should be fine since it's private to begin with)
