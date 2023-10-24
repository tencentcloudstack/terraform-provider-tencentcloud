---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_project"
sidebar_current: "docs-tencentcloud-datasource-bi_project"
description: |-
  Use this data source to query detailed information of bi project
---

# tencentcloud_bi_project

Use this data source to query detailed information of bi project

## Example Usage

```hcl
data "tencentcloud_bi_project" "project" {
  page_no           = 1
  keyword           = "abc"
  all_page          = true
  module_collection = "sys_common_user"
}
```

## Argument Reference

The following arguments are supported:

* `all_page` - (Optional, Bool) Whether to display all, if true, ignore paging.
* `keyword` - (Optional, String) Retrieve fuzzy fields.
* `module_collection` - (Optional, String) Role information, can be ignored.
* `page_no` - (Optional, Int) Page number.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `extra` - Additional information(Note: This field may return null, indicating that no valid value can be obtained).
* `list` - Array(Note: This field may return null, indicating that no valid value can be obtained).
  * `apply` - Apply(Note: This field may return null, indicating that no valid value can be obtained).
  * `auth_list` - List of permissions within the project(Note: This field may return null, indicating that no valid value can be obtained).
  * `color_code` - Logo colour(Note: This field may return null, indicating that no valid value can be obtained).
  * `config_list` - Customized parameters, this parameter can be ignored(Note: This field may return null, indicating that no valid value can be obtained).
    * `components` - Components(Note: This field may return null, indicating that no valid value can be obtained).
      * `include_type` - Include type(Note: This field may return null, indicating that no valid value can be obtained).
      * `module_id` - Module id(Note: This field may return null, indicating that no valid value can be obtained).
      * `params` - Extra parameters(Note: This field may return null, indicating that no valid value can be obtained).
    * `module_group` - Module group(Note: This field may return null, indicating that no valid value can be obtained).
  * `corp_id` - Enterprise id(Note: This field may return null, indicating that no valid value can be obtained).
  * `created_at` - Created at(Note: This field may return null, indicating that no valid value can be obtained).
  * `created_user` - Created by(Note: This field may return null, indicating that no valid value can be obtained).
  * `id` - Project id.
  * `is_external_manage` - Determine whether it is hosted(Note: This field may return null, indicating that no valid value can be obtained).
  * `last_modify_name` - Last modified report and presentation names(Note: This field may return null, indicating that no valid value can be obtained).
  * `logo` - Project logo(Note: This field may return null, indicating that no valid value can be obtained).
  * `manage_platform` - Hosting platform name(Note: This field may return null, indicating that no valid value can be obtained).
  * `mark` - Remark(Note: This field may return null, indicating that no valid value can be obtained).
  * `member_count` - Member count(Note: This field may return null, indicating that no valid value can be obtained).
  * `name` - Project name(Note: This field may return null, indicating that no valid value can be obtained).
  * `page_count` - Page count(Note: This field may return null, indicating that no valid value can be obtained).
  * `panel_scope` - Default kanban(Note: This field may return null, indicating that no valid value can be obtained).
  * `seed` - Obfuscated field(Note: This field may return null, indicating that no valid value can be obtained).
  * `source` - Interface call source(Note: This field may return null, indicating that no valid value can be obtained).
  * `updated_at` - Updated by(Note: This field may return null, indicating that no valid value can be obtained).
  * `updated_user` - Updated by(Note: This field may return null, indicating that no valid value can be obtained).
* `msg` - Interface information(Note: This field may return null, indicating that no valid value can be obtained).


