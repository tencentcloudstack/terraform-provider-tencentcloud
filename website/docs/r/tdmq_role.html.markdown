---
subcategory: "TDMQ for Pulsar(tpulsar)"
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
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
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

