Provides a resource to create a cfw sync_route

Example Usage

```hcl
resource "tencentcloud_cfw_sync_route" "example" {
  sync_type = "Route"
  fw_type   = "nat"
}
```