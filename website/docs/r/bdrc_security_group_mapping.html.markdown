---
subcategory: "Backup and Disaster Recovery Center(BDRC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_security_group_mapping"
sidebar_current: "docs-tencentcloud-resource-bdrc_security_group_mapping"
description: |-
  Provides a resource to create a BDRC security group mapping.
---

# tencentcloud_bdrc_security_group_mapping

Provides a resource to create a BDRC security group mapping.

## Example Usage

```hcl
resource "tencentcloud_bdrc_security_group_mapping" "example" {
  site_pair_id             = "sitepair-a4mtozsz"
  src_security_group_id    = "sg-ool7tmf8"
  target_security_group_id = "sg-jfy3gi92"
}
```

## Argument Reference

The following arguments are supported:

* `site_pair_id` - (Required, String, ForceNew) Security group mapping belongs to the site pair ID.
* `src_security_group_id` - (Required, String, ForceNew) Production end instance bound to the security group ID.
* `target_security_group_id` - (Required, String, ForceNew) Disaster recovery end instance bound to the security group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `life_state` - Life state of the security group mapping; NORMAL: normal.
* `security_group_mapping_id` - Security group mapping ID.
* `source_security_group_id` - Production end security group ID.


## Import

BDRC security group mapping can be imported using the sitePairId#securityGroupMappingId, e.g.

```
terraform import tencentcloud_bdrc_security_group_mapping.example sitepair-a4mtozsz#sgmap-88ylio5h
```

