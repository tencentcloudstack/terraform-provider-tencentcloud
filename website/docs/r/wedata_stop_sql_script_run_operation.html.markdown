---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_stop_sql_script_run_operation"
sidebar_current: "docs-tencentcloud-resource-wedata_stop_sql_script_run_operation"
description: |-
  Provides a resource to create a WeData stop sql script run operation
---

# tencentcloud_wedata_stop_sql_script_run_operation

Provides a resource to create a WeData stop sql script run operation

## Example Usage

```hcl
resource "tencentcloud_wedata_stop_sql_script_run_operation" "example" {
  job_id     = "ac13aceb-7a30-4414-91c0-6504f177462f"
  project_id = "2983848457986924544"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Specifies the query id.
* `project_id` - (Required, String, ForceNew) Project ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



