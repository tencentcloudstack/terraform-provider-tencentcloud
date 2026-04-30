Provides a resource to create TEO prefetch cache task.

Example Usage

```hcl
resource "tencentcloud_teo_prefetch_task_operation" "example" {
  zone_id = "zone-12345678"
  targets = [
    "http://www.example.com/example.txt",
  ]
}
```

Prefetch with edge mode and headers

```hcl
resource "tencentcloud_teo_prefetch_task_operation" "example" {
  zone_id = "zone-12345678"
  targets = [
    "http://www.example.com/example.txt",
  ]
  mode    = "edge"
  headers {
    name  = "X-Custom-Header"
    value = "custom-value"
  }
  prefetch_media_segments = "on"
}
```
