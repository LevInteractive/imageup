## API

There are only two endpoints.

* [POST](#process-file)
* [DELETE](#remove-files)

### Process File

Send a `POST` to `/` with the following arguments.

**Arguments:**

* `file`<string> - Binary file data.
* `sizes`<array> - Array of file configurations to return.
  * `width`<int> - Required. Max width for image.
  * `height`<int> - Required. Max height for image.
  * `fill`<bool> - If true, resize and crop to width and height using the [Lanczos resampling filter](https://github.com/disintegration/imaging#image-resizing).
  * `name`<string> - As no effect, but is returned with payload if provided for convenience.

```json
[
  {
    "width": 1000,
    "height": 1000,
    "fill": false
  },
  {
    "width": 500,
    "height": 500,
    "fill": true
  }
]
```

At the moment, processing options are minimal based on our needs, but more could
be added. Please [submit a issue for PR for more features](https://github.com/LevInteractive/imageup/issues).

**Example by cURL:**

```shell
curl -X POST \
  http://localhost:31111/ \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F file=@_DSC1814.jpg \
  -F 'sizes=[{"name": "foobar", "width":100, "height": 100, "fill": false}]'
```

## Response

The response objects will always be in the same order they were sent. The
`fileName` property contains the new generated unqiue name for this file. This
is the same value that should be stored and used to remove the file froms
storage.

```json
[
  {
    "fileName": "generatedname.jpg",
    "url": "https://the-public-url",
    "width": 1000,
    "height": 1000,
    "fill": true
  },
  {
    "fileName": "generatedname.jpg",
    "url": "https://the-public-url",
    "width": 500,
    "height": 500,
    "fill": true
  }
]
```
### Remove files

Send a `DELETE` to `/` with the following argument.

**Arguments:**

* `files` - A comma-delimited list of files. "file1, file2, etc"

**Example by cURL:**

```shell
curl -X DELETE \
  http://localhost:31111/ \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'files=foobar-82f100f-4bfd-4671-a0c1-9fc662782.jpg,foobar-38ef100f-4bfd-4671-a0c1-9fc667f2d924.jpg'
```

This will purge both files from the bucket.


## Errors

The `code` will be a relevant [http code](https://golang.org/pkg/net/http/#pkg-constants).

```json
{
    "code": 405,
    "message": "invalid"
}
```
