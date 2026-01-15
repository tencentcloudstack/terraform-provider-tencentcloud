Provides a resource to create a PostgreSQL parameter template

Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "PRC"
  }

  modify_param_entry_set {
    name           = "lock_timeout"
    expected_value = "60"
  }

  modify_param_entry_set {
    name           = "event_triggers"
    expected_value = "on"
  }
}
```

If you want delete some param like `lock_timeout`

```hcl
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "PRC"
  }

  modify_param_entry_set {
    name           = "event_triggers"
    expected_value = "on"
  }

  delete_param_set = ["lock_timeout"]
}
```

Import

PostgreSQL parameter template can be imported using the id, e.g.

Notice: `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.example 81ec47ed-0e4e-5af2-a648-2072fe63f225
```
