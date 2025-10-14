---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_ops_stop_task_async"
sidebar_current: "docs-tencentcloud-resource-wedata_ops_stop_task_async"
description: |-
  Provides a resource to create a wedata ops stop task async
---

# tencentcloud_wedata_ops_stop_task_async

Provides a resource to create a wedata ops stop task async

## Example Usage

```hcl
resource "tencentcloud_wedata_ops_stop_task_async" "wedata_ops_stop_task_async" {
  project_id = "1859317240494305280"
  task_ids   = ["20251013154418424"]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project id.
* `task_ids` - (Required, Set: [`String`], ForceNew) Task id list.
* `kill_instance` - (Optional, Bool, ForceNew) Whether to terminate the generated instance, the default is false; if true, it will wait for all forces to terminate.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



