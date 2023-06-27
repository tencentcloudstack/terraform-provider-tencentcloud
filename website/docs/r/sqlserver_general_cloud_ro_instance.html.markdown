---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_cloud_ro_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_cloud_ro_instance"
description: |-
  Provides a resource to create a sqlserver general_cloud_ro_instance
---

# tencentcloud_sqlserver_general_cloud_ro_instance

Provides a resource to create a sqlserver general_cloud_ro_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id                      = ""
  zone                             = ""
  read_only_group_type             =
  memory                           =
  storage                          =
  cpu                              =
  machine_type                     = ""
  read_only_group_forced_upgrade   =
  read_only_group_id               = ""
  read_only_group_name             = ""
  read_only_group_is_offline_delay =
  read_only_group_max_delay_time   =
  read_only_group_min_in_group     =
  instance_charge_type             = ""
  goods_num                        =
  subnet_id                        = ""
  vpc_id                           = ""
  period                           =
  security_group_list              =
  auto_voucher                     =
  voucher_ids                      =
  resource_tags {
    tag_key   = ""
    tag_value = ""

  }
  collation = ""
  time_zone = ""
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) Number of instance cores.
* `instance_id` - (Required, String) Primary instance ID, in the format: mssql-3l3fgqn7.
* `machine_type` - (Required, String) The host disk type of the purchased instance, CLOUD_HSSD-enhanced SSD cloud disk for virtual machines, CLOUD_TSSD-extremely fast SSD cloud disk for virtual machines, CLOUD_BSSD-universal SSD cloud disk for virtual machines.
* `memory` - (Required, Int) Instance memory size, in GB.
* `read_only_group_type` - (Required, Int) Read-only group type option, 1- Ship according to one instance and one read-only group, 2- Ship after creating a read-only group, all instances are under this read-only group, 3- All instances shipped are in the existing Some read-only groups below.
* `storage` - (Required, Int) Instance disk size, in GB.
* `zone` - (Required, String) Instance Availability Zone, similar to ap-guangzhou-1 (Guangzhou District 1); the instance sales area can be obtained through the interface DescribeZones.
* `collation` - (Optional, String) System character set collation, default: Chinese_PRC_CI_AS.
* `instance_charge_type` - (Optional, String) Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).
* `period` - (Optional, Int) Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48.
* `read_only_group_forced_upgrade` - (Optional, Int) 0 - Default not to upgrade the master instance, 1 - Mandatory upgrade of the master instance to complete ro deployment; if the master instance is a non-cluster version, you need to fill in 1 to force the upgrade to a cluster version. Filling in 1 indicates that you have agreed to upgrade the master instance to a cluster instance.
* `read_only_group_id` - (Optional, String) Required when ReadOnlyGroupType=3, existing read-only group ID.
* `read_only_group_is_offline_delay` - (Optional, Int) Required when ReadOnlyGroupType=2, whether to enable the delayed elimination function for the newly created read-only group, 1-on, 0-off. When the delay between the read-only replica and the primary instance is greater than the threshold, it will be automatically removed.
* `read_only_group_max_delay_time` - (Optional, Int) Mandatory when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the threshold for delay culling of newly created read-only groups.
* `read_only_group_min_in_group` - (Optional, Int) Required when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the newly created read-only group retains at least the number of read-only replicas after delay elimination.
* `read_only_group_name` - (Optional, String) Required when ReadOnlyGroupType=2, the name of the newly created read-only group.
* `resource_tags` - (Optional, List) A collection of tags bound to the new instance.
* `security_group_list` - (Optional, Set: [`String`]) Security group list, fill in the security group ID in the form of sg-xxx.
* `subnet_id` - (Optional, String) VPC subnet ID, in the form of subnet-bdoe83fa; SubnetId and VpcId need to be set at the same time or not set at the same time.
* `time_zone` - (Optional, String) System time zone, default: China Standard Time.
* `vpc_id` - (Optional, String) VPC network ID, in the form of vpc-dsp338hz; SubnetId and VpcId need to be set at the same time or not set at the same time.

The `resource_tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



