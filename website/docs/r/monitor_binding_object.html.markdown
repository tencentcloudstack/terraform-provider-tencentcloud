---
subcategory: "Monitor"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_binding_object"
sidebar_current: "docs-tencentcloud-resource-monitor_binding_object"
description: |-
  Provides a resource for bind objects to a policy group resource.
---

# tencentcloud_monitor_binding_object

Provides a resource for bind objects to a policy group resource.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_monitor_policy_binding_object.

## Example Usage

```hcl
data "tencentcloud_instances" "instances" {
}
resource "tencentcloud_monitor_policy_group" "group" {
  group_name       = "terraform_test"
  policy_view_name = "cvm_device"
  remark           = "this is a test policy group"
  is_union_rule    = 1
  conditions {
    metric_id           = 33
    alarm_notify_type   = 1
    alarm_notify_period = 600
    calc_type           = 1
    calc_value          = 3
    calc_period         = 300
    continue_period     = 2
  }
}

#for cvm
resource "tencentcloud_monitor_binding_object" "binding" {
  group_id = tencentcloud_monitor_policy_group.group.id
  dimensions {
    dimensions_json = "{\"unInstanceId\":\"${data.tencentcloud_instances.instances.instance_list[0].instance_id}\"}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dimensions` - (Required, Set, ForceNew) A list objects. Each element contains the following attributes:
* `group_id` - (Required, Int, ForceNew) Policy group ID for binding objects.

The `dimensions` object supports the following:

* `dimensions_json` - (Required, String, ForceNew) Represents a collection of dimensions of an object instance, json format.eg:'{"unInstanceId":"ins-ot3cq4bi"}'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



