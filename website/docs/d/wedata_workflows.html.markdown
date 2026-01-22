---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_workflows"
sidebar_current: "docs-tencentcloud-datasource-wedata_workflows"
description: |-
  Use this data source to query detailed information of wedata wedata_workflows
---

# tencentcloud_wedata_workflows

Use this data source to query detailed information of wedata wedata_workflows

## Example Usage

```hcl
data "tencentcloud_wedata_workflows" "wedata_workflows" {
  project_id = 2905622749543821312
  keyword    = "test_workflow"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `bundle_id` - (Optional, String) bundleId item.
* `create_time` - (Optional, Set: [`String`]) Creation time range yyyy-MM-dd HH:MM:ss. two times must be filled in the array.
* `create_user_uin` - (Optional, String) Creator ID.
* `keyword` - (Optional, String) Search keywords.
* `modify_time` - (Optional, Set: [`String`]) Modification time interval yyyy-MM-dd HH:MM:ss. fill in two times in the array.
* `owner_uin` - (Optional, String) Owner ID.
* `parent_folder_path` - (Optional, String) Workflow folder.
* `result_output_file` - (Optional, String) Used to save results.
* `workflow_type` - (Optional, String) Workflow type. valid values: cycle and manual.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Describes workflow pagination information.


