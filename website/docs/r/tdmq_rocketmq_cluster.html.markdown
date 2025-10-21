---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_cluster"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_cluster"
description: |-
  Provides a resource to create a tdmqRocketmq cluster
---

# tencentcloud_tdmq_rocketmq_cluster

Provides a resource to create a tdmqRocketmq cluster

## Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) Cluster name, which can contain 3-64 letters, digits, hyphens, and underscores.
* `remark` - (Optional, String) Cluster description (up to 128 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cluster_id` - Cluster ID.
* `create_time` - Creation time in milliseconds.
* `is_vip` - Whether it is an exclusive instance.
* `public_end_point` - Public network access address.
* `region` - Region information.
* `rocket_m_q_flag` - Rocketmq cluster identification.
* `support_namespace_endpoint` - Whether the namespace access point is supported.
* `vpc_end_point` - VPC access address.
* `vpcs` - Vpc list.
  * `subnet_id` - Subnet ID.
  * `vpc_id` - Vpc ID.


## Import

tdmqRocketmq cluster can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_cluster.cluster cluster_id
```

