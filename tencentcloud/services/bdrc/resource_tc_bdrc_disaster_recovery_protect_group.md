Provides a resource to create a BDRC disaster recovery protect group.

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = "sitepair-xxxxxxxx"
  protect_group_type       = "DISK"
  recovery_point_objective = 15
  protect_group_name       = "tf-example-protect-group"
  data_direction           = "POSITIVE"
}
```

Import

BDRC disaster recovery protect group can be imported using the ProtectGroupId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_protect_group.example pg-xxxxxxxx
```
