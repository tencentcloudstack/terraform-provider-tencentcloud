Provides a resource to create a dnspod snapshot_config

Example Usage

```hcl
resource "tencentcloud_dnspod_snapshot_config" "snapshot_config" {
  domain = "dnspod.cn"
  period = "hourly"
}
```

Import

dnspod snapshot_config can be imported using the id, e.g.

```
terraform import tencentcloud_dnspod_snapshot_config.snapshot_config domain
```