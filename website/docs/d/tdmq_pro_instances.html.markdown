---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_pro_instances"
sidebar_current: "docs-tencentcloud-datasource-tdmq_pro_instances"
description: |-
  Use this data source to query detailed information of tdmq pro_instances
---

# tencentcloud_tdmq_pro_instances

Use this data source to query detailed information of tdmq pro_instances

## Example Usage

```hcl
data "tencentcloud_tdmq_pro_instances" "pro_instances_filter" {
  filters {
    name   = "InstanceName"
    values = ["keep"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) query condition filter.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) The name of the filter parameter.
* `values` - (Optional, Set) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - Instance information list.
  * `auto_renew_flag` - Automatic renewal mark, 0 indicates the default state (the user has not set it, that is, the initial state is manual renewal), 1 indicates automatic renewal, 2 indicates that the automatic renewal is not specified (user setting).
  * `config_display` - Instance configuration specification name.
  * `create_time` - Create time.
  * `expire_time` - Instance expiration time, in milliseconds.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_version` - Instance version.
  * `max_band_width` - Peak bandwidth. Unit: mbps.
  * `max_storage` - Storage capacity, in GB.
  * `max_tps` - Peak TPS.
  * `pay_mode` - 0-postpaid, 1-prepaid.
  * `remark` - RemarksNote: This field may return null, indicating that no valid value can be obtained.
  * `scalable_tps` - Elastic TPS outside specificationNote: This field may return null, indicating that no valid value can be obtained.
  * `spec_name` - Instance Configuration ID.
  * `status` - Instance status, 0-creating, 1-normal, 2-isolating, 3-destroyed, 4-abnormal, 5-delivery failure, 6-allocation change, 7-allocation failure.
  * `subnet_id` - Subnet idNote: This field may return null, indicating that no valid value can be obtained.
  * `tags` - Tag list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `vpc_id` - Id of the VPCNote: This field may return null, indicating that no valid value can be obtained.


