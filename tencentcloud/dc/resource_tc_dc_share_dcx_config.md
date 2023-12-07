Provides a resource to create a dc share_dcx_config

Example Usage

```hcl
resource "tencentcloud_dc_share_dcx_config" "share_dcx_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  enable = false
}
```

Import

dc share_dcx_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_share_dcx_config.share_dcx_config dcx_id
```