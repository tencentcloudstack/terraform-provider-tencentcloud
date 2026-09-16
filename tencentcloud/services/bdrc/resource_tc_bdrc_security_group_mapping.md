Provides a resource to create a BDRC security group mapping.

Example Usage

```hcl
resource "tencentcloud_bdrc_security_group_mapping" "example" {
  site_pair_id             = "sitepair-a4mtozsz"
  src_security_group_id    = "sg-ool7tmf8"
  target_security_group_id = "sg-jfy3gi92"
}
```

Import

BDRC security group mapping can be imported using the sitePairId#securityGroupMappingId, e.g.

```
terraform import tencentcloud_bdrc_security_group_mapping.example sitepair-a4mtozsz#sgmap-88ylio5h
```
