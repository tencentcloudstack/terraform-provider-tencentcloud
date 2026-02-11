Provides a resource to create a PostgreSQL parameter template config

~> **NOTE:** The `tencentcloud_postgresql_parameter_template_config` and `tencentcloud_postgresql_parameter_template` resources are mutually exclusive: if one is used to configure a parameter template, the other cannot be used simultaneously.

Example Usage

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

Import

PostgreSQL parameter template config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_parameter_template_config.example 0c595485-c1b8-518b-bd87-dfe44a530fa5
```
