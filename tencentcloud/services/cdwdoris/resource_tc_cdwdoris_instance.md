Provides a resource to create a cdwdoris cdwdoris_instance

Example Usage

```hcl
resource "tencentcloud_cdwdoris_instance" "cdwdoris_instance" {
  fe_spec = {
  }
  be_spec = {
  }
  charge_properties = {
  }
  tags = {
  }
  user_multi_zone_infos = {
  }
}
```

Import

cdwdoris cdwdoris_instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdwdoris_instance.cdwdoris_instance cdwdoris_instance_id
```
