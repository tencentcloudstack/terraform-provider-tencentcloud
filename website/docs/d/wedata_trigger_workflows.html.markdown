---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_trigger_workflows"
sidebar_current: "docs-tencentcloud-datasource-wedata_trigger_workflows"
description: |-
  Use this data source to query detailed information of wedata trigger workflows
---

# tencentcloud_wedata_trigger_workflows

Use this data source to query detailed information of wedata trigger workflows

## Example Usage

```hcl
data "tencentcloud_wedata_trigger_workflows" "wedata_trigger_workflows" {
  project_id = "3108707295180644352"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `bundle_id` - (Optional, String) Bundle ID item.
* `create_time` - (Optional, Set: [`String`]) Creation time range yyyy-MM-dd HH:mm:ss, two timestamps need to be filled in the array.
* `create_user_uin` - (Optional, String) Creator ID.
* `keyword` - (Optional, String) Search keyword.
* `modify_time` - (Optional, Set: [`String`]) Modification time range yyyy-MM-dd HH:mm:ss, two timestamps need to be filled in the array.
* `owner_uin` - (Optional, String) Owner ID.
* `parent_folder_path` - (Optional, String) Folder path to which the workflow belongs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Paginated workflow query information.


