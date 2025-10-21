---
subcategory: "Project"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_projects"
sidebar_current: "docs-tencentcloud-datasource-projects"
description: |-
  Use this data source to query detailed information of tag project
---

# tencentcloud_projects

Use this data source to query detailed information of tag project

## Example Usage

```hcl
data "tencentcloud_projects" "project" {
  all_list = 1
}
```

## Argument Reference

The following arguments are supported:

* `all_list` - (Required, Int) 1 means to list all project, 0 means to list visible project.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `projects` - List of projects.
  * `create_time` - Create time.
  * `creator_uin` - Uin of Creator.
  * `project_id` - ID of Project.
  * `project_info` - Description of project.
  * `project_name` - Name of Project.


