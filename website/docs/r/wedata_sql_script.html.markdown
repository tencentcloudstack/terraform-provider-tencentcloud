---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_sql_script"
sidebar_current: "docs-tencentcloud-resource-wedata_sql_script"
description: |-
  Provides a resource to create a WeData sql script
---

# tencentcloud_wedata_sql_script

Provides a resource to create a WeData sql script

## Example Usage

```hcl
resource "tencentcloud_wedata_sql_folder" "example" {
  folder_name        = "tf_example"
  project_id         = "2983848457986924544"
  parent_folder_path = "/"
  access_scope       = "SHARED"
}

resource "tencentcloud_wedata_sql_script" "example" {
  script_name        = "tf_example_script"
  project_id         = "2983848457986924544"
  parent_folder_path = tencentcloud_wedata_sql_folder.example.path
  script_config {
    datasource_id    = "108826"
    compute_resource = "svmgao_stability"
  }

  script_content = "SHOW DATABASES;"
  access_scope   = "SHARED"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String, ForceNew) Project ID.
* `script_name` - (Required, String, ForceNew) Script name.
* `access_scope` - (Optional, String, ForceNew) Permission scope: SHARED, PRIVATE.
* `parent_folder_path` - (Optional, String, ForceNew) Parent folder path, /aaa/bbb/ccc, root directory is empty string or /.
* `script_config` - (Optional, List) Data exploration script configuration.
* `script_content` - (Optional, String) Script content, if there is a value.

The `script_config` object supports the following:

* `advance_config` - (Optional, String) Advanced settings, execution configuration parameters, map-json String,String. Encoded in Base64.
* `compute_resource` - (Optional, String) Computing resource.
* `datasource_env` - (Optional, String) Data source environment.
* `datasource_id` - (Optional, String) Data source ID.
* `executor_group_id` - (Optional, String) Execution resource group.
* `params` - (Optional, String) Advanced runtime parameters, variable substitution, map-json String,String.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `path` - The full path of the node, /aaa/bbb/ccc.ipynb, consists of the names of each node.
* `script_id` - Script ID.


## Import

WeData sql script can be imported using the projectId#scriptId, e.g.

```
terraform import tencentcloud_wedata_sql_script.example 2983848457986924544#cccc3170-6334-46c3-adce-c5776ad2280c
```

