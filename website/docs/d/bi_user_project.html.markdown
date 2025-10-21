---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_user_project"
sidebar_current: "docs-tencentcloud-datasource-bi_user_project"
description: |-
  Use this data source to query detailed information of bi user_project
---

# tencentcloud_bi_user_project

Use this data source to query detailed information of bi user_project

## Example Usage

```hcl
data "tencentcloud_bi_user_project" "user_project" {
  project_id = 123
  all_page   = true
}
```

## Argument Reference

The following arguments are supported:

* `all_page` - (Optional, Bool) Whether to display all, if true, ignore paging.
* `project_id` - (Optional, Int) Project id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Array(Note: This field may return null, indicating that no valid value can be obtained).
  * `area_code` - Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).
  * `corp_id` - Enterprise id(Note: This field may return null, indicating that no valid value can be obtained).
  * `created_at` - Created at(Note: This field may return null, indicating that no valid value can be obtained).
  * `created_user` - Created by(Note: This field may return null, indicating that no valid value can be obtained).
  * `email` - E-mail(Note: This field may return null, indicating that no valid value can be obtained).
  * `first_modify` - First login to change password, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).
  * `global_user_name` - Global role name(Note: This field may return null, indicating that no valid value can be obtained).
  * `last_login` - Last login time, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).
  * `mobile` - Mobile number, public cloud unrelated fields(Note: This field may return null, indicating that no valid value can be obtained).
  * `phone_number` - Phone number(Note: This field may return null, indicating that no valid value can be obtained).
  * `status` - Disabled state(Note: This field may return null, indicating that no valid value can be obtained).
  * `updated_at` - Updated at(Note: This field may return null, indicating that no valid value can be obtained).
  * `updated_user` - Updated by(Note: This field may return null, indicating that no valid value can be obtained).
  * `user_id` - User id.
  * `user_name` - Username.


