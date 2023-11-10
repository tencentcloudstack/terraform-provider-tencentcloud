---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_restart_push_task"
sidebar_current: "docs-tencentcloud-resource-css_restart_push_task"
description: |-
  Provides a resource to create a css restart_push_task
---

# tencentcloud_css_restart_push_task

Provides a resource to create a css restart_push_task

## Example Usage

```hcl
resource "tencentcloud_css_restart_push_task" "restart_push_task" {
  task_id  = "3d5738dd-1ca2-4601-a6e9-004c5ec75c0b"
  operator = "tf-test"
}
```

## Argument Reference

The following arguments are supported:

* `operator` - (Required, String, ForceNew) Task operator.
* `task_id` - (Required, String, ForceNew) Task Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



