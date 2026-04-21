Provides a resource to check CNAME status for TEO domains.

Example Usage

```hcl
resource "tencentcloud_teo_check_cname_status_operation" "example" {
  zone_id = "zone-12345678"
  record_names = [
    "example.com",
    "www.example.com",
  ]
}
```
