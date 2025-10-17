---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_run_sql_script_operation"
sidebar_current: "docs-tencentcloud-resource-wedata_run_sql_script_operation"
description: |-
  Provides a resource to create a WeData run sql script operation
---

# tencentcloud_wedata_run_sql_script_operation

Provides a resource to create a WeData run sql script operation

## Example Usage

```hcl
resource "tencentcloud_wedata_run_sql_script_operation" "example" {
  script_id  = "195a5f49-8e43-4e05-8b42-cecdfb6349f8"
  project_id = "2983848457986924544"
}
```

### Or

```hcl
resource "tencentcloud_wedata_run_sql_script_operation" "example" {
  script_id      = "195a5f49-8e43-4e05-8b42-cecdfb6349f8"
  project_id     = "2983848457986924544"
  script_content = "SHOW DATABASES;"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `script_id` - (Required, String, ForceNew) Script id.
* `params` - (Optional, String, ForceNew) Advanced running parameter.
* `script_content` - (Optional, String, ForceNew) Script content. executed by default if not transmitted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Job ID of the SQL script operation.


