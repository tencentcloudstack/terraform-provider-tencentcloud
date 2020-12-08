---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_attachment"
description: |-
  Provides a resource to create a CLB attachment.
---

# tencentcloud_clb_attachment

Provides a resource to create a CLB attachment.

## Example Usage

```hcl
resource "tencentcloud_clb_attachment" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-hh141sn9"
  rule_id     = "loc-4xxr2cy7"

  targets {
    instance_id = "ins-1flbqyp8"
    port        = 80
    weight      = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of the CLB.
* `listener_id` - (Required, ForceNew) Id of the CLB listener.
* `targets` - (Required) Information of the backends to be attached.
* `rule_id` - (Optional, ForceNew) Id of the CLB listener rule. Only supports listeners of 'HTTPS' and 'HTTP' protocol.

The `targets` object supports the following:

* `instance_id` - (Required) Id of the backend server.
* `port` - (Required) Port of the backend server. Valid value ranges: (0~65535).
* `weight` - (Optional) Forwarding weight of the backend service. Valid value ranges: (0~100). defaults to 10.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `protocol_type` - Type of protocol within the listener.


## Import

CLB attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_attachment.foo loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
```

