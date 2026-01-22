---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_task_set_success_instance_async"
sidebar_current: "docs-tencentcloud-resource-wedata_task_set_success_instance_async"
description: |-
  Provides a resource to create a wedata task set success instance
---

# tencentcloud_wedata_task_set_success_instance_async

Provides a resource to create a wedata task set success instance

## Example Usage

```hcl
resource "tencentcloud_wedata_task_set_success_instance_async" "wedata_task_set_success_instance_async" {
  project_id        = "1859317240494305280"
  instance_key_list = ["20250324192240178_2025-10-13 17:00:00"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_key_list` - (Required, Set: [`String`], ForceNew) Instance id list, which can be obtained from ListInstances.
* `project_id` - (Required, String, ForceNew) Project Id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



