---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_submit_trigger_task"
sidebar_current: "docs-tencentcloud-resource-wedata_submit_trigger_task"
description: |-
  Provides a resource to create a wedata wedata_submit_trigger_task
---

# tencentcloud_wedata_submit_trigger_task

Provides a resource to create a wedata wedata_submit_trigger_task

## Example Usage

```hcl
resource "tencentcloud_wedata_submit_trigger_task" "wedata_submit_trigger_task" {
  project_id     = "2905622749543821312"
  task_id        = "20251015164958429"
  version_remark = "v1"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `task_id` - (Required, String, ForceNew) Task ID.
* `version_remark` - (Required, String, ForceNew) Version remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



