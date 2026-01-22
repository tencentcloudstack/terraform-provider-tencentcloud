Provides a resource to create a cynosdb backup_config

Example Usage

Enable logical backup configuration and cross-region logical backup

```hcl
resource "tencentcloud_cynosdb_backup_config" "foo" {
    backup_time_beg  = 7200
    backup_time_end  = 21600
    cluster_id       = "cynosdbmysql-bws8h88b"
    reserve_duration = 604800

    logic_backup_config {
        logic_backup_enable        = "ON"
        logic_backup_time_beg      = 7200
        logic_backup_time_end      = 21600
        logic_cross_regions        = ["ap-shanghai"]
        logic_cross_regions_enable = "ON"
        logic_reserve_duration     = 259200
    }
}
```

Disable logical backup configuration

```hcl
resource "tencentcloud_cynosdb_backup_config" "foo" {
    backup_time_beg  = 7200
    backup_time_end  = 21600
    cluster_id       = "cynosdbmysql-bws8h88b"
    reserve_duration = 604800

    logic_backup_config {
        logic_backup_enable        = "OFF"
    }
}
```

Import

cynosdb backup_config can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_backup_config.foo cynosdbmysql-bws8h88b
```