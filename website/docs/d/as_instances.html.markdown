---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_instances"
sidebar_current: "docs-tencentcloud-datasource-as_instances"
description: |-
  Use this data source to query detailed information of as instances
---

# tencentcloud_as_instances

Use this data source to query detailed information of as instances

## Example Usage

```hcl
data "tencentcloud_as_instances" "instances" {
  filters {
    name   = "auto-scaling-group-id"
    values = [tencentcloud_as_scaling_group.scaling_group.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions. If there are multiple Filters, the relationship between Filters is a logical AND (AND) relationship. If there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR (OR) relationship.
* `instance_ids` - (Optional, Set: [`String`]) Instance ID of the cloud server (CVM) to be queried. The limit is 100 per request.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Fields to be filtered. Valid names: `instance-id`: Filters by instance ID, `auto-scaling-group-id`: Filter by scaling group ID.
* `values` - (Required, Set) Value of the field.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - List of instance details.
  * `add_time` - The time when the instance joined the group.
  * `auto_scaling_group_id` - Auto scaling group ID.
  * `auto_scaling_group_name` - Auto scaling group name.
  * `creation_type` - Valid values: `AUTO_CREATION`, `MANUAL_ATTACHING`.
  * `health_status` - Health status, the valid values are HEALTHY and UNHEALTHY.
  * `instance_id` - Instance ID.
  * `instance_type` - Instance type.
  * `launch_configuration_id` - Launch configuration ID.
  * `launch_configuration_name` - Launch configuration name.
  * `life_cycle_state` - Life cycle state. Please refer to the link for field value details: https://cloud.tencent.com/document/api/377/20453#Instance.
  * `protected_from_scale_in` - Enable scale in protection.
  * `version_number` - Version ID.
  * `zone` - Available zone.


