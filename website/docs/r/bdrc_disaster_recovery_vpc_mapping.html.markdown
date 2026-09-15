---
subcategory: "Backup & Disaster Recovery(BDRC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_disaster_recovery_vpc_mapping"
sidebar_current: "docs-tencentcloud-resource-bdrc_disaster_recovery_vpc_mapping"
description: |-
  Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC.
---

# tencentcloud_bdrc_disaster_recovery_vpc_mapping

Provides a resource to create a disaster recovery VPC mapping for site pair of BDRC.

## Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_vpc_mapping" "example" {
  site_pair_id     = "site-pair-xxxx"
  source_vpc_id    = "vpc-source-xxxx"
  source_subnet_id = "subnet-source-xxxx"
  target_vpc_id    = "vpc-target-xxxx"
  target_subnet_id = "subnet-target-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `site_pair_id` - (Required, String, ForceNew) Site pair ID of the disaster recovery VPC mapping.
* `source_subnet_id` - (Required, String, ForceNew) Source subnet ID of the disaster recovery VPC mapping.
* `source_vpc_id` - (Required, String, ForceNew) Source VPC ID of the disaster recovery VPC mapping.
* `target_subnet_id` - (Required, String, ForceNew) Target subnet ID of the disaster recovery VPC mapping.
* `target_vpc_id` - (Required, String, ForceNew) Target VPC ID of the disaster recovery VPC mapping.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - Mapping rule primary key ID.
* `life_state` - Lifecycle state.
* `source_subnet` - Source subnet ID returned by the cloud API.
* `source_vpc` - Source VPC ID returned by the cloud API.
* `status` - Mapping status.
* `target_subnet` - Target subnet ID returned by the cloud API.
* `target_vpc` - Target VPC ID returned by the cloud API.


## Import

BDRC disaster recovery VPC mapping can be imported using the compound id sitePairId#vpcMappingId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_vpc_mapping.example site-pair-xxxx#88
```

