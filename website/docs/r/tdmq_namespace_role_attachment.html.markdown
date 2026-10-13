---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_namespace_role_attachment"
sidebar_current: "docs-tencentcloud-resource-tdmq_namespace_role_attachment"
description: |-
  Provide a resource to create a TDMQ environment role.
---

# tencentcloud_tdmq_namespace_role_attachment

Provide a resource to create a TDMQ environment role.

## Example Usage

```hcl
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
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}

resource "tencentcloud_tdmq_namespace_role_attachment" "example" {
  environ_id  = tencentcloud_tdmq_namespace.example.environ_name
  role_name   = tencentcloud_tdmq_role.example.role_name
  permissions = ["produce", "consume"]
  cluster_id  = tencentcloud_tdmq_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The id of tdmq cluster.
* `environ_id` - (Required, String) The name of tdmq namespace.
* `permissions` - (Required, List: [`String`]) The permissions of tdmq role.
* `role_name` - (Required, String) The name of tdmq role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of resource.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEnvironmentRole](https://cloud.tencent.com/document/api/1179/62402) - 创建环境角色授权
* Read: [DescribeEnvironmentRoles](https://cloud.tencent.com/document/api/1179/49051) - 获取命名空间角色列表
* Update: [ModifyEnvironmentRole](https://cloud.tencent.com/document/api/1179/62398) - 修改环境角色授权
* Delete: [DeleteEnvironmentRoles](https://cloud.tencent.com/document/api/1179/62400) - 删除环境角色授权

