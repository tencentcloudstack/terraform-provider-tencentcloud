---
subcategory: "CLB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_target_group_attachment"
description: |-
  Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.
---

# tencentcloud_clb_target_group_attachment

Provides a resource to create a CLB target group attachment is bound to the load balancing listener or forwarding rule.

~> **NOTE:** Required argument `targrt_group_id` is no longer supported, replace by `target_group_id`.

## Example Usage

```hcl
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-rule-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.listener_basic.id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}

resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test-target-keep-1"
}

resource "tencentcloud_clb_target_group_attachment" "group" {
  clb_id          = tencentcloud_clb_instance.clb_basic.id
  listener_id     = tencentcloud_clb_listener.listener_basic.id
  rule_id         = tencentcloud_clb_listener_rule.rule_basic.id
  target_group_id = tencentcloud_clb_target_group.test.id
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Required, ForceNew) ID of the CLB.
* `listener_id` - (Required, ForceNew) ID of the CLB listener.
* `target_group_id` - (Required, ForceNew) ID of the CLB target group.
* `rule_id` - (Optional, ForceNew) ID of the CLB listener rule.
* `targrt_group_id` - (Optional, ForceNew, **Deprecated**) It has been deprecated from version 1.47.1. Use `target_group_id` instead. ID of the CLB target group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB target group attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_attachment.group lbtg-odareyb2#lbl-bicjmx3i#lb-cv0iz74c#loc-ac6uk7b6
```

