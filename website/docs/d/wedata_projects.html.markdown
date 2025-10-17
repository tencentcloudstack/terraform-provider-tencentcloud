---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_projects"
sidebar_current: "docs-tencentcloud-datasource-wedata_projects"
description: |-
  Use this data source to query detailed information of WeData projects
---

# tencentcloud_wedata_projects

Use this data source to query detailed information of WeData projects

## Example Usage

### Query all projects

```hcl
data "tencentcloud_wedata_projects" "example" {}
```

### Query projects by filter

```hcl
data "tencentcloud_wedata_projects" "example" {
  project_ids = [
    "2982667120655491072",
    "2853989879663501312"
  ]

  project_name  = "tf_example"
  status        = 1
  project_model = "SIMPLE"
}
```

## Argument Reference

The following arguments are supported:

* `project_ids` - (Optional, Set: [`String`]) List of project IDs.
* `project_model` - (Optional, String) Project model, optional values: SIMPLE, STANDARD.
* `project_name` - (Optional, String) Project name or unique identifier name, supports fuzzy search.
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Int) Project status, optional values: 0 (disabled), 1 (normal).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - List of data sources.
  * `create_time` - Creation time.
  * `creator_uin` - Project creator ID.
  * `description` - Remarks.
  * `display_name` - Project display name, can be Chinese name.
  * `project_id` - Project ID.
  * `project_model` - Project model, SIMPLE: simple mode, STANDARD: standard mode.
  * `project_name` - Project identifier, English name.
  * `project_owner_uin` - Project owner ID.
  * `status` - Project status: 0: disabled, 1: enabled, -3: disabling, 2: enabling.


