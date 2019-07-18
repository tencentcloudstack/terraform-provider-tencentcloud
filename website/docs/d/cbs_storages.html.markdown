---
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

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The available zone that the CBS instance locates at.
* `project_id` - (Optional) ID of the project with which the CBS is associated.
* `result_output_file` - (Optional) Used to save results.
* `storage_id` - (Optional) ID of the CBS to be queried.
* `storage_name` - (Optional) Name of the CBS to be queried.
* `storage_type` - (Optional) Types of storage medium, and available values include CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.
* `storage_usage` - (Optional) Types of CBS, and available values include SYSTEM_DISK and DATA_DISK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `storage_list` - A list of storage. Each element contains the following attributes:
  * `attached` - Indicates whether the CBS is mounted the CVM.
  * `availability_zone` - The zone of CBS.
  * `create_time` - Creation time of CBS.
  * `encrypt` - Indicates whether CBS is encrypted.
  * `instance_id` - ID of the CVM instance that be mounted by this CBS.
  * `project_id` - ID of the project.
  * `status` - Status of CBS.
  * `storage_id` - ID of CBS.
  * `storage_name` - Name of CBS.
  * `storage_size` - Volume of CBS.
  * `storage_type` - Types of storage medium.
  * `storage_usage` - Types of CBS.
  * `tags` - The available tags within this CBS.


