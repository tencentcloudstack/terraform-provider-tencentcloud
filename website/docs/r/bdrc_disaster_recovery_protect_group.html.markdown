---
subcategory: "Backup and Disaster Recovery Center(BDRC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_disaster_recovery_protect_group"
sidebar_current: "docs-tencentcloud-resource-bdrc_disaster_recovery_protect_group"
description: |-
  Provides a resource to create a BDRC disaster recovery protect group.
---

# tencentcloud_bdrc_disaster_recovery_protect_group

Provides a resource to create a BDRC disaster recovery protect group.

## Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = "sitepair-a4mtozsz"
  protect_group_name       = "tf-example"
  data_direction           = "POSITIVE"
  protect_group_type       = "INSTANCE"
  recovery_point_objective = 15
}
```

## Argument Reference

The following arguments are supported:

* `protect_group_type` - (Required, String, ForceNew) Product type of the disaster recovery protect group. Valid values: `DISK`, `INSTANCE`, `CFS`.
* `recovery_point_objective` - (Required, Int, ForceNew) Expected RPO of the protect group, in minutes (currently only 15 is supported).
* `site_pair_id` - (Required, String, ForceNew) The ID of the disaster recovery site pair to which the protect group belongs.
* `data_direction` - (Optional, String, ForceNew) Data replication direction. Valid values: `POSITIVE`, `REVERSE`.
* `protect_group_name` - (Optional, String) Name of the protect group, up to 60 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_uin` - Account Uin of the protect group owner.
* `app_id` - User AppId.
* `bind_protected_resource_count` - Number of protected resources bound to the protect group.
* `copy_type` - Replication technology (SYN sync / ASY async).
* `create_from` - Creation source (LOCAL local / PEER peer).
* `create_time` - Creation time.
* `disaster_recovery_type` - Disaster recovery type (CROSS_ZONE / CROSS_REGION / CROSS_CLOUD).
* `error_recovery_point_objective_count` - Number of replication pairs whose RPO is abnormal (not synced for more than 15 minutes).
* `life_state` - Lifecycle state.
* `modify_time` - Modification time.
* `peer_cloud_name` - Peer cloud name (only returned when DisasterRecoveryType is CROSS_CLOUD).
* `protected_resource_status_set` - Protected resource status statistics, key is the replication pair status, value is the resource count under that status.
  * `count` - Resource count under this status.
  * `status` - Replication pair status.
* `site_pair_name` - Name of the disaster recovery site pair.
* `source_region` - Source region.
* `source_vpc` - Source VPC.
* `source_zone` - Source zone.
* `sub_account_uin` - Sub account Uin of the protect group creator.
* `target_region` - Target region.
* `target_vpc` - Target VPC.
* `target_zone` - Target zone.


## Import

BDRC disaster recovery protect group can be imported using the protectGroupId#protectGroupType, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_protect_group.example pg-l5xwdgsn#INSTANCE
```

