---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_password_complexity"
sidebar_current: "docs-tencentcloud-resource-mysql_password_complexity"
description: |-
  Provides a resource to create a mysql password_complexity
---

# tencentcloud_mysql_password_complexity

Provides a resource to create a mysql password_complexity

## Example Usage

```hcl
resource "tencentcloud_mysql_password_complexity" "password_complexity" {
  instance_id = var.instance_id
  param_list {
    name          = "validate_password_length"
    current_value = "8"
  }
  param_list {
    name          = "validate_password_mixed_case_count"
    current_value = "2"
  }
  param_list {
    name          = "validate_password_number_count"
    current_value = "2"
  }
  param_list {
    name          = "validate_password_special_char_count"
    current_value = "2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `param_list` - (Optional, List) List of parameters to be modified. Every element is a combination of `Name` (parameter name) and `CurrentValue` (new value). Valid values for `Name` of version 8.0: `validate_password.policy`, `validate_password.lengt`, `validate_password.mixed_case_coun`, `validate_password.number_coun`, `validate_password.special_char_count`. Valid values for `Name` of version 5.6 and 5.7: `validate_password_polic`, `validate_password_lengt` `validate_password_mixed_case_coun`, `validate_password_number_coun`, `validate_password_special_char_coun`.

The `param_list` object supports the following:

* `current_value` - (Optional, String) Parameter value.
* `name` - (Optional, String) Parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



