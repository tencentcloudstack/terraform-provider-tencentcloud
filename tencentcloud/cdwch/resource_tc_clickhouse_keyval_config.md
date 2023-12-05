Provides a resource to create a clickhouse keyval_config

Example Usage

```hcl
resource "tencentcloud_clickhouse_keyval_config" "keyval_config" {
  instance_id = "cdwch-datuhk3z"
  items {
    conf_key   = "max_open_files"
    conf_value = "50000"
  }
}
```

Import

clickhouse config can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_keyval_config.config cdwch-datuhk3z#max_open_files#50000
```