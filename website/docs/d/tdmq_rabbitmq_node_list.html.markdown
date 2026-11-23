---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_node_list"
sidebar_current: "docs-tencentcloud-datasource-tdmq_rabbitmq_node_list"
description: |-
  Use this data source to query detailed information of tdmq rabbitmq_node_list
---

# tencentcloud_tdmq_rabbitmq_node_list

Use this data source to query detailed information of tdmq rabbitmq_node_list

## Example Usage

```hcl
data "tencentcloud_tdmq_rabbitmq_node_list" "rabbitmq_node_list" {
  instance_id = "amqp-testtesttest"
  node_name   = "keep-node"
  filters {
    name   = "nodeStatus"
    values = ["running", "down"]
  }
  sort_element = "cpuUsage"
  sort_order   = "descend"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) rabbitmq cluster ID.
* `filters` - (Optional, List) filter parameter name and valueNow there is only one nodeStatusrunning/downArray type, compatible with adding filter parameters later.
* `node_name` - (Optional, String) Fuzzy search node name.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_element` - (Optional, String) Sort by the specified element, now there are only 2cpuUsage/diskUsage.
* `sort_order` - (Optional, String) ascend/descend.

The `filters` object supports the following:

* `name` - (Optional, String) The name of the filter parameter.
* `values` - (Optional, Set) value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_list` - cluster listNote: This field may return null, indicating that no valid value can be obtained.
  * `cpu_usage` - CPU usageNote: This field may return null, indicating that no valid value can be obtained.
  * `disk_usage` - disk usageNote: This field may return null, indicating that no valid value can be obtained.
  * `memory` - Memory usage, in MBNote: This field may return null, indicating that no valid value can be obtained.
  * `node_name` - node nameNote: This field may return null, indicating that no valid value can be obtained.
  * `node_status` - node statusNote: This field may return null, indicating that no valid value can be obtained.
  * `process_number` - Number of Erlang processes for RabbitmqNote: This field may return null, indicating that no valid value can be obtained.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeRabbitMQNodeList](https://www.tencentcloud.com/search/?wd=DescribeRabbitMQNodeList&type=document&lang=en)

