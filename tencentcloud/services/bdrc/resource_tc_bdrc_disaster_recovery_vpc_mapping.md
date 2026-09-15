Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC.

Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_vpc_mapping" "example" {
  site_pair_id     = "site-pair-xxxx"
  source_vpc_id    = "vpc-source-xxxx"
  source_subnet_id = "subnet-source-xxxx"
  target_vpc_id    = "vpc-target-xxxx"
  target_subnet_id = "subnet-target-xxxx"
}
```

Import

BDRC disaster recovery VPC mapping can be imported using the compound id sitePairId#vpcMappingId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_vpc_mapping.example site-pair-xxxx#88
```
