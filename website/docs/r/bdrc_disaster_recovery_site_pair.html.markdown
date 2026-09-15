---
subcategory: "Backup and Disaster Recovery Center(BDRC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_disaster_recovery_site_pair"
sidebar_current: "docs-tencentcloud-resource-bdrc_disaster_recovery_site_pair"
description: |-
  Provides a resource to create a BDRC (Backup and Disaster Recovery Center) disaster recovery site pair.
---

# tencentcloud_bdrc_disaster_recovery_site_pair

Provides a resource to create a BDRC (Backup and Disaster Recovery Center) disaster recovery site pair.

## Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_site_pair" "example" {
  disaster_recovery_type = "CROSS_REGION"
  source_region          = "ap-guangzhou"
  source_zone            = "ap-guangzhou-3"
  target_region          = "ap-shanghai"
  target_zone            = "ap-shanghai-2"
  source_vpc             = "vpc-xxxxxxxx"
  target_vpc             = "vpc-yyyyyyyy"
  site_pair_product_type = "DISK"
  site_pair_name         = "tf-example-site-pair"
  copy_type              = "ASY"
}
```

## Argument Reference

The following arguments are supported:

* `disaster_recovery_type` - (Required, String, ForceNew) Disaster recovery type, CROSS_REGION (cross-region) or CROSS_ZONE (cross-zone).
* `site_pair_product_type` - (Required, String, ForceNew) Site pair product type, including DISK, CFS, INSTANCE.
* `source_region` - (Required, String, ForceNew) Production site region.
* `source_vpc` - (Required, String, ForceNew) Production site VPC.
* `source_zone` - (Required, String, ForceNew) Production site availability zone.
* `target_region` - (Required, String, ForceNew) Disaster recovery site region.
* `target_vpc` - (Required, String, ForceNew) Disaster recovery site VPC.
* `target_zone` - (Required, String, ForceNew) Disaster recovery site availability zone.
* `copy_type` - (Optional, String, ForceNew) Replication technology, SYN (synchronous) / ASY (asynchronous).
* `site_pair_name` - (Optional, String) Site pair name, max length is 60 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_uin` - Master account Uin of the account that created the site pair.
* `bind_protect_group_count` - Number of bound protection groups.
* `create_from` - Creation source, LOCAL (local creation) / PEER (peer creation).
* `create_time` - Creation time.
* `cross_cloud_details` - Extra information for cross-cloud scenarios (only returned when IsCrossCloud=true).
  * `peer_region_name` - Peer cloud region display name.
  * `peer_vpc_name` - Peer cloud VPC display name.
  * `peer_zone_name` - Peer cloud availability zone display name.
  * `source_app_id` - Source cloud AppId.
  * `source_cloud_name` - Source cloud name (peer cloud name in cross-cloud).
  * `source_sub_account_uin` - Source cloud sub account Uin.
  * `source_uin` - Source cloud master account Uin.
  * `source_user_name` - Source cloud user name.
  * `target_app_id` - Target cloud AppId.
  * `target_cloud_name` - Target cloud name (local cloud name in cross-cloud).
  * `target_sub_account_uin` - Target cloud sub account Uin.
  * `target_uin` - Target cloud master account Uin.
* `error_recovery_point_objective_copy_pair_set` - List of copy pair IDs with RPO anomalies.
* `protected_resource_set` - List of protected resources grouped by resource type.
  * `resource_id_set` - List of protected source resource IDs under this type.
  * `resource_type` - Resource type (consistent with SitePairType, such as DISK/CFS/INSTANCE).
* `protected_resource_status_set` - Status statistics of protected resources.
  * `count` - Number of resources under this status.
  * `status` - Copy pair status.
* `site_pair_state` - Site pair state.
* `site_pair_type` - Site pair type (product type, such as DISK/CFS/INSTANCE).
* `sub_account_uin` - Sub account Uin of the account that created the site pair.


## Import

BDRC disaster recovery site pair can be imported using the SitePairId, e.g.

```
terraform import tencentcloud_bdrc_disaster_recovery_site_pair.example site-pair-xxxxxx
```

