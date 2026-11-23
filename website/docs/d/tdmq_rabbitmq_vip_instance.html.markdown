---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_vip_instance"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rabbitmq_vip_instance"
description: |-
  Use this data source to query detailed information of tdmq rabbitmq_vip_instance
---

# tencentcloud_tdmq_rabbitmq_vip_instance

Use this data source to query detailed information of tdmq rabbitmq_vip_instance

## Example Usage

```hcl
data "tencentcloud_tdmq_rabbitmq_vip_instance" "rabbitmq_vip_instance" {
  filters {
    name   = ""
    values = []
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
  * `exception_information` - The cluster is abnormal.Note: This field may return null, indicating that no valid value can be obtained.
  * `expire_time` - Instance expiration time, in milliseconds.
  * `instance_id` - instance id.
  * `instance_name` - instance name.
  * `instance_version` - instance versionNote: This field may return null, indicating that no valid value can be obtained.
  * `max_band_width` - Peak bandwidth, in Mbps.
  * `max_storage` - Storage capacity, in GB.
  * `max_tps` - Peak TPS.
  * `node_count` - Number of nodes.
  * `pay_mode` - 0-postpaid, 1-prepaid.
  * `remark` - RemarksNote: This field may return null, indicating that no valid value can be obtained.
  * `spec_name` - Instance Configuration ID.
  * `status` - Instance status, 0 means creating, 1 means normal, 2 means isolating, 3 means destroyed, 4 - abnormal, 5 - delivery failed.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeRabbitMQVipInstances](https://www.tencentcloud.com/search/?wd=DescribeRabbitMQVipInstances&type=document&lang=en)

