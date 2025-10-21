---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_projects"
sidebar_current: "docs-tencentcloud-datasource-dcdb_projects"
description: |-
  Use this data source to query detailed information of dcdb projects
---

# tencentcloud_dcdb_projects

Use this data source to query detailed information of dcdb projects

## Example Usage

```hcl
data "tencentcloud_dcdb_projects" "projects" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `projects` - Project list.
  * `app_id` - Application ID.
  * `create_time` - Creation time.
  * `creator_uin` - Creator UIN.
  * `info` - Description.
  * `is_default` - Whether it is the default project. Valid values: `1` (yes), `0` (no).
  * `name` - Project name.
  * `owner_uin` - The UIN of the resource owner (root account).
  * `project_id` - Project ID.
  * `src_app_id` - Source APPID.
  * `src_plat` - Source platform.
  * `status` - Project status. Valid values: `0` (normal), `-1` (disabled), `3` (default project).


