Provides a resource to create a postgresql parameter_template

Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name = "your_temp_name"
  db_major_version = "13"
  db_engine = "postgresql"
  template_description = "For_tf_test"

  modify_param_entry_set {
	name = "timezone"
	expected_value = "UTC"
  }
  modify_param_entry_set {
	name = "lock_timeout"
	expected_value = "123"
  }

  delete_param_set = ["lc_time"]
}
```

Import

postgresql parameter_template can be imported using the id, e.g.

Notice: `modify_param_entry_set` and `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.parameter_template parameter_template_id
```