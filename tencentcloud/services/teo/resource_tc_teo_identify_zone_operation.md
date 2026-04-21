Provides a resource to identify TEO zone or subdomain ownership.

Example Usage

```hcl
resource "tencentcloud_teo_identify_zone_operation" "example" {
  zone_name = "example.com"
}
```

With subdomain

```hcl
resource "tencentcloud_teo_identify_zone_operation" "example_sub" {
  zone_name = "example.com"
  domain    = "www.example.com"
}
```
