---
subcategory: "TDMQ"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_instance"
sidebar_current: "docs-tencentcloud-resource-tdmq_instance"
description: |-
  Provide a resource to create a TDMQ instance.
---

# tencentcloud_tdmq_instance

Provide a resource to create a TDMQ instance.

## Example Usage

```hcl
resource "tencentcloud_tdmq_instance" "foo" {
  cluster_name = "example"
  remark       = "this is description."
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) The name of tdmq cluster to be created.
* `bind_cluster_id` - (Optional, Int) The Dedicated Cluster Id.
* `remark` - (Optional, String) Description of the tdmq cluster.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Tdmq instance can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test tdmq_id
```

