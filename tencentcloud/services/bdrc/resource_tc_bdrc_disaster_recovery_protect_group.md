Provides a resource to create a BDRC disaster recovery protect group.

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = "sitepair-a4mtozsz"
  protect_group_name       = "tf-example"
  data_direction           = "POSITIVE"
  protect_group_type       = "INSTANCE"
  recovery_point_objective = 15
}
```

Import

BDRC disaster recovery protect group can be imported using the protectGroupId#protectGroupType, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_protect_group.example pg-l5xwdgsn#INSTANCE
```
