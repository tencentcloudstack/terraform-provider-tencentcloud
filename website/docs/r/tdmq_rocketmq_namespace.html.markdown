---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rocketmq_namespace"
sidebar_current: "docs-tencentcloud-resource-tdmq_rocketmq_namespace"
description: |-
  Provides a resource to create a tdmqRocketmq namespace
---

# tencentcloud_tdmq_rocketmq_namespace

Provides a resource to create a tdmqRocketmq namespace

## Example Usage

```hcl
resource "tencentcloud_tdmq_rocketmq_cluster" "cluster" {
  cluster_name = "test_rocketmq"
  remark       = "test recket mq"
}

resource "tencentcloud_tdmq_rocketmq_namespace" "namespace" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.cluster.cluster_id
  namespace_name = "test_namespace"
  ttl            = 65000
  retention_time = 65000
  remark         = "test namespace"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace_name` - (Required, String) Namespace name, which can contain 3-64 letters, digits, hyphens, and underscores.
* `retention_time` - (Required, Int) Retention time of persisted messages in milliseconds.
* `ttl` - (Required, Int) Retention time of unconsumed messages in milliseconds. Value range: 60 seconds-15 days.
* `remark` - (Optional, String) Remarks (up to 128 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_endpoint` - Public network access point address.
* `vpc_endpoint` - VPC access point address.


## Import

tdmqRocketmq namespace can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_namespace.namespace namespace_id
```

