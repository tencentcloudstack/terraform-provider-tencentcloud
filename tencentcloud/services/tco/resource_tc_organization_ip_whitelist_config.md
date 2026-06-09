Provides a resource to create a Organization IP whitelist config

Example Usage

```hcl
resource "tencentcloud_organization_ip_whitelist_config" "example" {
  zone_id = "z-1os7c9znogct"
  ip_whitelist = [
    "10.0.0.0/24",
    "192.168.1.0/24",
    "172.16.10.0/24",
  ]
}
```

Import

Organization IP whitelist config can be imported using the zoneId, e.g.

```
terraform import tencentcloud_organization_ip_whitelist_config.example z-1os7c9znogct
```