---
subcategory: "Tencent Distributed Message Queue(TDMQ)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_role"
sidebar_current: "docs-tencentcloud-resource-tdmq_role"
description: |-
  Provide a resource to create a TDMQ role.
---

# tencentcloud_tdmq_role

Provide a resource to create a TDMQ role.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark       = "this is description."
}

resource "tencentcloud_tdmq_namespace" "bar" {
  environ_name = "example"
  msg_ttl      = 300
  cluster_id   = "${tencentcloud_tdmq_instance.foo.id}"
  remark       = "this is description."
}

resource "tencentcloud_tdmq_topic" "bar" {
  environ_id = "${tencentcloud_tdmq_namespace.bar.id}"
  topic_name = "example"
  partitions = 6
  topic_type = 0
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark     = "this is description."
}

resource "tencentcloud_tdmq_role" "bar" {
  role_name  = "example"
  cluster_id = "${tencentcloud_tdmq_instance.foo.id}"
  remark     = "this is description world"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The id of tdmq cluster.
* `remark` - (Required, String) The description of tdmq role.
* `role_name` - (Required, String) The name of tdmq role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test tdmq_id
```

