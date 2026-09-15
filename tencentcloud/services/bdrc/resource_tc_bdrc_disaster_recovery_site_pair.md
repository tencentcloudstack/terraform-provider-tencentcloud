Provides a resource to create a BDRC (Backup and Disaster Recovery Center) disaster recovery site pair.

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_site_pair" "example" {
  disaster_recovery_type   = "CROSS_REGION"
  source_region            = "ap-guangzhou"
  source_zone              = "ap-guangzhou-3"
  target_region            = "ap-shanghai"
  target_zone              = "ap-shanghai-2"
  source_vpc               = "vpc-xxxxxxxx"
  target_vpc               = "vpc-yyyyyyyy"
  site_pair_product_type   = "DISK"
  site_pair_name           = "tf-example-site-pair"
  copy_type                = "ASY"
}
```

Import

BDRC disaster recovery site pair can be imported using the SitePairId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_site_pair.example site-pair-xxxxxx
```
