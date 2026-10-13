---
subcategory: "TDSQL-C for PostgreSQL(TDCPG)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdcpg_instance"
sidebar_current: "docs-tencentcloud-resource-tdcpg_instance"
description: |-
  Provides a resource to create a tdcpg instance.
---

# tencentcloud_tdcpg_instance

Provides a resource to create a tdcpg instance.

~> **NOTE:** This resource is still in internal testing. To experience its functions, you need to apply for a whitelist from Tencent Cloud.

## Example Usage

```hcl
resource "tencentcloud_tdcpg_instance" "instance1" {
  cluster_id    = "cluster_id"
  cpu           = 1
  memory        = 1
  instance_name = "instance_name"
}

resource "tencentcloud_tdcpg_instance" "instance2" {
  cluster_id       = "cluster_id"
  cpu              = 1
  memory           = 2
  instance_name    = "instance_name"
  operation_timing = "IMMEDIATE"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) cluster id.
* `cpu` - (Required, Int) cpu cores.
* `memory` - (Required, Int) memory size.
* `instance_name` - (Optional, String) instance name.
* `operation_timing` - (Optional, String) operation timing, optional value is IMMEDIATE or MAINTAIN_PERIOD.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateClusterInstances](https://cloud.tencent.com/document/api/1556/70690) - 创建实例
* Create: [DescribeResourcesByDealName](https://cloud.tencent.com/document/api/1556/70679) - 根据订单号获取资源信息
* Read: [DescribeClusterInstances](https://cloud.tencent.com/document/api/1556/70688) - 查询实例
* Update: [DescribeClusterInstances](https://cloud.tencent.com/document/api/1556/70688) - 查询实例
* Update: [ModifyClusterInstancesSpec](https://cloud.tencent.com/document/api/1556/70686) - 修改实例规格
* Delete: [DeleteClusterInstances](https://cloud.tencent.com/document/api/1556/70689) - 删除实例
* Delete: [DescribeClusterInstances](https://cloud.tencent.com/document/api/1556/70688) - 查询实例
* Delete: [IsolateClusterInstances](https://cloud.tencent.com/document/api/1556/70687) - 隔离实例

## Import

tdcpg instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_tdcpg_instance.instance cluster_id#instance_id
```

