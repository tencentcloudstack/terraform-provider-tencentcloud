Provides a resource to create a BDRC security group mapping.

Example Usage

```hcl
resource "tencentcloud_bdrc_security_group_mapping" "example" {
  site_pair_id           = "sp-xxxxxx"
  src_security_group_id  = "sg-xxxxxxxx"
  target_security_group_id = "sg-yyyyyyyy"
}
```

Import

BDRC security group mapping can be imported using the sitePairId#securityGroupMappingId, e.g.

```
terraform import tencentcloud_bdrc_security_group_mapping.example sp-xxxxxx#sgm-yyyyyy
```
