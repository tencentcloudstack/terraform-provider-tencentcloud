Provides a resource to create a BDRC disaster recovery protect group.

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_site_pair" "example" {
  disaster_recovery_type = "CROSS_ZONE"
  source_region          = "ap-shanghai"
  source_zone            = "ap-shanghai-3"
  target_region          = "ap-shanghai"
  target_zone            = "ap-shanghai-4"
  source_vpc             = "vpc-lx6q09ji"
  target_vpc             = "vpc-jktad5e6"
  site_pair_product_type = "INSTANCE"
  site_pair_name         = "tf-example"
  copy_type              = "ASY"
}

resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
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
