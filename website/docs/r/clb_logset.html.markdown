---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_logset"
sidebar_current: "docs-tencentcloud-resource-clb_logset"
description: |-
  Provides a resource to create an exclusive CLB Logset.
---

# tencentcloud_clb_logset

Provides a resource to create an exclusive CLB Logset.

## Example Usage

```hcl
resource "tencentcloud_clb_logset" "foo" {
  name    = "clb_logset"
  perioid = 7
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Logset name, which must be unique among all CLS logsets. Default is `clb_logset`.
* `period` - (Optional, ForceNew) Logset retention period in days. Maximun value is `90`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Logset creation time.
* `topic_count` - Number of log topics in logset.


## Import

CLB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_logset.foo 4eb9e3a8-9c42-4b32-9ddf-e215e9c92764
```

