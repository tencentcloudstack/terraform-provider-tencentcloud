Provides a resource to create TEO cache purge task.

Example Usage

Purge URLs

```hcl
resource "tencentcloud_teo_purge_task" "purge_url_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_url"
  targets = [
    "https://example.com/path1",
    "https://example.com/path2",
  ]
}
```

Purge all cache

```hcl
resource "tencentcloud_teo_purge_task" "purge_all_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_all"
}
```

Purge cache tag

```hcl
resource "tencentcloud_teo_purge_task" "purge_cache_tag_example" {
  zone_id = "zone-2qtuhspy7cr6"
  type    = "purge_cache_tag"
  cache_tag {
    domains = [
      "example.com",
      "www.example.com",
    ]
  }
}
```
