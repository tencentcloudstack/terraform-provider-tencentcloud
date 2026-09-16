Provides a resource to create a BDRC (Backup and Disaster Recovery Center) disaster recovery site pair.

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
```

Import

BDRC disaster recovery site pair can be imported using the sitePairId#sitePairProductType, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_site_pair.example sitepair-la2re1mv#INSTANCE
```
