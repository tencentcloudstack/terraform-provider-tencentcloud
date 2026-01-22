Provides a resource to create a dc share dcx config

Example Usage

```hcl
resource "tencentcloud_dc_share_dcx_config" "example" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  enable                   = true
}
```

Import

dc share dcx config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_share_dcx_config.example dcx-4z49tnws
```