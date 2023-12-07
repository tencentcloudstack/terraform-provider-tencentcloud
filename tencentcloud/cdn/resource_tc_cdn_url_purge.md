Provide a resource to invoke a Url Purge Request.

Example Usage

```hcl
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
}
```

Change `redo` argument to request new purge task with same urls

```hcl
resource "tencentcloud_cdn_url_purge" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
  redo = 1
}
```