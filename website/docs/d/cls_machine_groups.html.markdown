---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_machine_groups"
sidebar_current: "docs-tencentcloud-datasource-cls_machine_groups"
description: |-
  Use this data source to query detailed information of cls machine_groups
---

# tencentcloud_cls_machine_groups

Use this data source to query detailed information of cls machine_groups

## Example Usage

### Query all machine groups

```hcl
data "tencentcloud_cls_machine_groups" "example" {}
```

### Query machine group by filters

```hcl
data "tencentcloud_cls_machine_groups" "example" {
  filters {
    name   = "machineGroupId"
    values = ["76e09d6f-e0c5-4103-bd6e-22bdbf63a76e"]
  }

  filters {
    name   = "machineGroupName"
    values = ["cls-m26ybna6"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. Maximum 10 filters, each with up to 5 values. Supported keys: `machineGroupName`, `machineGroupId`, `osType`, `tagKey`, `tag:tagKey`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name.
* `values` - (Required, Set) Filter field values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `machine_groups` - List of cls machine groups.
  * `auto_update` - Whether machine group auto update is enabled.
  * `create_time` - Creation time.
  * `delay_cleanup_time` - Machine offline periodic cleanup time, in days.
  * `group_id` - Machine group ID.
  * `group_name` - Machine group name.
  * `machine_group_type` - Machine group type.
    * `type` - Machine group type. Valid values: `ip`, `label`.
    * `values` - Machine description list.
  * `meta_tags` - Machine group metadata tag list.
    * `key` - Metadata tag key.
    * `value` - Metadata tag value.
  * `os_type` - Operating system type. 0: Linux, 1: Windows.
  * `service_logging` - Whether service logging is enabled.
  * `tags` - Tag list.
    * `key` - Tag key.
    * `value` - Tag value.
  * `update_end_time` - Upgrade end time.
  * `update_start_time` - Upgrade start time.


