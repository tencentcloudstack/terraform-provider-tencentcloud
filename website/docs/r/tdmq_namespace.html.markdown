---
subcategory: "TDMQ for Pulsar(tpulsar)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_namespace"
sidebar_current: "docs-tencentcloud-resource-tdmq_namespace"
description: |-
  Provide a resource to create a tdmq namespace.
---

# tencentcloud_tdmq_namespace

Provide a resource to create a tdmq namespace.

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
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The Dedicated Cluster Id.
* `environ_name` - (Required, String) The name of namespace to be created.
* `msg_ttl` - (Required, Int) The expiration time of unconsumed message.
* `remark` - (Optional, String) Description of the namespace.
* `retention_policy` - (Optional, Map) The Policy of message to retain.

The `retention_policy` object supports the following:

* `size_in_mb` - (Optional, Int) the size of message to retain.
* `time_in_minutes` - (Optional, Int) the time of message to retain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Tdmq namespace can be imported, e.g.

```
$ terraform import tencentcloud_tdmq_instance.test namespace_id
```

