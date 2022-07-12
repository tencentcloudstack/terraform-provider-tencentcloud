---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storages"
sidebar_current: "docs-tencentcloud-datasource-cbs_storages"
description: |-
  Use this data source to query detailed information of CBS storages.
---

# tencentcloud_cbs_storages

Use this data source to query detailed information of CBS storages.

## Example Usage

```hcl
data "tencentcloud_cbs_storages" "storages" {
  storage_id         = "disk-kdt0sq6m"
  result_output_file = "mytestpath"
}
```

The following snippet shows the new supported query params

```hcl
data "tencentcloud_cbs_storages" "whats_new" {
  charge_type   = ["POSTPAID_BY_HOUR", "PREPAID"]
  portable      = true
  storage_state = ["ATTACHED"]
  instance_ips  = ["10.0.0.2"]
  instance_name = ["my-instance"]
  tag_keys      = ["foo"]
  tag_values    = ["bar", "baz"]
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the CBS instance locates at.
* `charge_type` - (Optional, List: [`String`]) List filter by disk charge type (`POSTPAID_BY_HOUR` | `PREPAID`).
* `instance_ips` - (Optional, List: [`String`]) List filter by attached instance public or private IPs.
* `instance_name` - (Optional, List: [`String`]) List filter by attached instance name.
* `portable` - (Optional, Bool) Filter by whether the disk is portable (Boolean `true` or `false`).
* `project_id` - (Optional, Int) ID of the project with which the CBS is associated.
* `result_output_file` - (Optional, String) Used to save results.
* `storage_id` - (Optional, String) ID of the CBS to be queried.
* `storage_name` - (Optional, String) Name of the CBS to be queried.
* `storage_state` - (Optional, List: [`String`]) List filter by disk state (`UNATTACHED` | `ATTACHING` | `ATTACHED` | `DETACHING` | `EXPANDING` | `ROLLBACKING` | `TORECYCLE`).
* `storage_type` - (Optional, String) Filter by cloud disk media type (`CLOUD_BASIC`: HDD cloud disk | `CLOUD_PREMIUM`: Premium Cloud Storage | `CLOUD_SSD`: SSD cloud disk).
* `storage_usage` - (Optional, String) Filter by cloud disk type (`SYSTEM_DISK`: system disk | `DATA_DISK`: data disk).
* `tag_keys` - (Optional, List: [`String`]) List filter by tag keys.
* `tag_values` - (Optional, List: [`String`]) List filter by tag values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `storage_list` - A list of storage. Each element contains the following attributes:
  * `attached` - Indicates whether the CBS is mounted the CVM.
  * `availability_zone` - The zone of CBS.
  * `charge_type` - Pay type of the CBS instance.
  * `create_time` - Creation time of CBS.
  * `encrypt` - Indicates whether CBS is encrypted.
  * `instance_id` - ID of the CVM instance that be mounted by this CBS.
  * `prepaid_renew_flag` - The way that CBS instance will be renew automatically or not when it reach the end of the prepaid tenancy.
  * `project_id` - ID of the project.
  * `status` - Status of CBS.
  * `storage_id` - ID of CBS.
  * `storage_name` - Name of CBS.
  * `storage_size` - Volume of CBS.
  * `storage_type` - Types of storage medium.
  * `storage_usage` - Types of CBS.
  * `tags` - The available tags within this CBS.
  * `throughput_performance` - Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.


