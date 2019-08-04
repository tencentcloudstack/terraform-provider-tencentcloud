---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_server_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_server_attachment"
description: |-
  Provide a resource to create a CLB instance.
---

# tencentcloud_clb_server_attachment

Provide a resource to create a CLB instance.

## Example Usage

```hcl
resource "tencentcloud_clb_server_attachment" "attachment" {
  listener_id   = "lbl-hh141sn9#lb-k2zjp9lv"
  clb_id        = "lb-k2zjp9lv"
  protocol_type = "tcp"
  location_id   = "loc-4xxr2cy7"
  targets = {
    instance_id = "ins-1flbqyp8"
    port        = 50
    weight      = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) Id of the cloud load balancer. 
* `listener_id` - (Required, ForceNew) Id of the cloud load balance listener. 
* `targets` - (Required) Backend infos.
* `location_id` - (Optional, ForceNew) Id of the cloud load balance listener rule. 

The `targets` object supports the following:

* `instance_id` - (Required) Id of the backend server.
* `port` - (Required) Port of the backend server.
* `weight` - (Optional) Weight of the backend server.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `protocol_type` - Type of protocol within the listener, and available values include 'TCP', 'UDP', 'HTTP', 'HTTPS' and 'TCP_SSL'. 


## Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_server_attachment.attachment loc-4xxr2cy7#lbl-hh141sn9#lb-7a0t6zqb
```

