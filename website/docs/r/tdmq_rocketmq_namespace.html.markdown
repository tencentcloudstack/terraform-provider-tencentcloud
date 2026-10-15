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
resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}

resource "tencentcloud_tdmq_rocketmq_namespace" "example" {
  cluster_id     = tencentcloud_tdmq_rocketmq_cluster.example.cluster_id
  namespace_name = "tf_example_namespace"
  remark         = "remark."
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `namespace_name` - (Required, String) Namespace name, which can contain 3-64 letters, digits, hyphens, and underscores.
* `remark` - (Optional, String) Remarks (up to 128 characters).
* `retention_time` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.81.20. Due to the adjustment of RocketMQ, the creation or modification of this parameter will be ignored. Retention time of persisted messages in milliseconds.
* `ttl` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.81.20. Due to the adjustment of RocketMQ, the creation or modification of this parameter will be ignored. Retention time of unconsumed messages in milliseconds. Value range: 60 seconds-15 days.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_endpoint` - Public network access point address.
* `vpc_endpoint` - VPC access point address.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRocketMQNamespace](https://cloud.tencent.com/document/api/1179/63427) - 创建 RocketMQ 命名空间
* Read: [DescribeRocketMQNamespaces](https://cloud.tencent.com/document/api/1179/63419) - 查询 RocketMQ 命名空间列表
* Update: [ModifyRocketMQNamespace](https://cloud.tencent.com/document/api/1179/63415) - 修改 RocketMQ 命名空间
* Delete: [DeleteRocketMQNamespace](https://cloud.tencent.com/document/api/1179/63412) - 删除 RocketMQ 命名空间

## Import

tdmqRocketmq namespace can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdmq_rocketmq_namespace.namespace namespace_id
```

