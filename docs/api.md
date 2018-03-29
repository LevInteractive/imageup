## API

### Add files to storage

Send a `POST` to `/` with the following arguments.

**Arguments:**

* `file` - Binary file data)
* `sizes` - Array of file configurations to return

The `file` will be converted to all of the dimensions specified in the `sizes`
argument. If the array below was sent, two images would be saved and returned.

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

You may include a `name` property which won't effect anything but can
be used for convinience as it will be passed along to the returned objects.

**Example by cURL:**

```shell
curl -X POST \
  http://localhost:31111/ \
  -H 'cache-control: no-cache' \
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
### Remove file(s) from storage

Send a `DELETE` to `/` with the following argument.

**Arguments:**

* `files` - A comma-delimited list of files. "file1, file2, etc"

**Example by cURL:**

```shell
curl -X DELETE \
  http://localhost:31111/ \
  -H 'cache-control: no-cache' \
  -H 'content-type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW' \
  -F 'files=foobar-82f100f-4bfd-4671-a0c1-9fc662782.jpg,foobar-38ef100f-4bfd-4671-a0c1-9fc667f2d924.jpg'
```


## Errors

The `code` will be a relevant [http code](https://golang.org/pkg/net/http/#pkg-constants).

```json
{
    "code": 405,
    "message": "invalid"
}
```
