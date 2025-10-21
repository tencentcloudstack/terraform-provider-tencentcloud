---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pull_stream_task_restart"
sidebar_current: "docs-tencentcloud-resource-css_pull_stream_task_restart"
description: |-
  Provides a resource to create a css restart_push_task
---

# tencentcloud_css_pull_stream_task_restart

Provides a resource to create a css restart_push_task

## Example Usage

```hcl
resource "tencentcloud_css_pull_stream_task_restart" "restart_push_task" {
  task_id  = "3573"
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



