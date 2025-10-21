---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_environment_attributes"
sidebar_current: "docs-tencentcloud-datasource-tdmq_environment_attributes"
description: |-
  Use this data source to query detailed information of tdmq environment_attributes
---

# tencentcloud_tdmq_environment_attributes

Use this data source to query detailed information of tdmq environment_attributes

## Example Usage

```hcl
data "tencentcloud_tdmq_environment_attributes" "example" {
  environment_id = tencentcloud_tdmq_namespace.example.environ_name
  cluster_id     = tencentcloud_tdmq_instance.example.id
}

resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  remark       = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `environment_id` - (Required, String) Environment (namespace) name.
* `cluster_id` - (Optional, String) ID of the Pulsar cluster.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `msg_ttl` - Expiration time of unconsumed messages, unit second, maximum 1296000 (15 days).
* `rate_in_byte` - Consumption rate limit, unit byte/second, 0 unlimited rate.
* `rate_in_size` - Consumption rate limit, unit number/second, 0 is unlimited.
* `remark` - Remark.
* `replicas` - Duplicate number.
* `retention_hours` - Consumed message storage policy, unit hour, 0 will be deleted immediately after consumption.
* `retention_size` - Consumed message storage strategy, unit G, 0 Delete immediately after consumption.


