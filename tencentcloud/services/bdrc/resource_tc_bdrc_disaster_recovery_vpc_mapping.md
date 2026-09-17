Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC.

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

resource "tencentcloud_bdrc_disaster_recovery_vpc_mapping" "example" {
  site_pair_id     = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  source_vpc_id    = "vpc-lx6q09ji"
  source_subnet_id = "subnet-nyrg9pkl"
  target_vpc_id    = "vpc-jktad5e6"
  target_subnet_id = "subnet-jdzuvvvb"
}
```

Import

BDRC disaster recovery VPC mapping can be imported using the compound id sitePairId#vpcMappingId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_vpc_mapping.example sitepair-1sxvs5oj#13
```
