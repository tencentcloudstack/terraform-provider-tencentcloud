---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_namespace_role_attachment"
sidebar_current: "docs-tencentcloud-resource-tdmq_namespace_role_attachment"
description: |-
  Provide a resource to create a TDMQ role.
---

# tencentcloud_tdmq_namespace_role_attachment

Provide a resource to create a TDMQ role.

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


