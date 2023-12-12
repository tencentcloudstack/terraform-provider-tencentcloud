Provide a resource to invoke a Url Push request.

Example Usage

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = ["https://www.example.com/b"]
}
```

Change `redo` argument to request new push task with same urls

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
  redo = 1
}
```