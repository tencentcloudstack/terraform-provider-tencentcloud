Provides a resource to create a dcdb db_sync_mode_config

Example Usage

```hcl
resource "tencentcloud_dcdb_db_sync_mode_config" "config" {
  instance_id = "%s"
  sync_mode = 2
}
```

Import

dcdb db_sync_mode_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_sync_mode_config.db_sync_mode_config db_sync_mode_config_id
```