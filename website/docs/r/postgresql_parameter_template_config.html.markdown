---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameter_template_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_parameter_template_config"
description: |-
  Provides a resource to create a PostgreSQL parameter template config
---

# tencentcloud_postgresql_parameter_template_config

Provides a resource to create a PostgreSQL parameter template config

~> **NOTE:** The `tencentcloud_postgresql_parameter_template_config` and `tencentcloud_postgresql_parameter_template` resources are mutually exclusive: if one is used to configure a parameter template, the other cannot be used simultaneously.

## Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."
}

resource "tencentcloud_postgresql_parameter_template_config" "example" {
  template_id = tencentcloud_postgresql_parameter_template.example.id
  modify_param_entry_set {
    name           = "min_parallel_index_scan_size"
    expected_value = "64"
  }

  modify_param_entry_set {
    name           = "enable_async_append"
    expected_value = "on"
  }

  modify_param_entry_set {
    name           = "enable_group_by_reordering"
    expected_value = "on"
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_id` - (Required, String, ForceNew) Specifies the parameter template ID, which uniquely identifies the parameter template and cannot be modified. it can be obtained through the api [DescribeParameterTemplates](https://www.tencentcloud.comom/document/api/409/84067?from_cn_redirect=1).
* `modify_param_entry_set` - (Optional, Set) The set of parameters to be modified or added.

The `modify_param_entry_set` object supports the following:

* `expected_value` - (Required, String) The new value to which the parameter will be modified. When this parameter is used as an input parameter, its value must be a string, such as `0.1` (decimal), `1000` (integer), and `replica` (enum).
* `name` - (Required, String) Parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

PostgreSQL parameter template config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_parameter_template_config.example 0c595485-c1b8-518b-bd87-dfe44a530fa5
```

