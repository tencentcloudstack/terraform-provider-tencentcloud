Provides a resource to create a postgresql backup_download_restriction_config

Example Usage

Unlimit the restriction of the backup file download.
```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "NONE"
}
```

Set the download only to allow the intranet downloads.
```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "INTRANET"
}
```

Restrict the backup file download by customizing.
```hcl
resource "tencentcloud_vpc" "pg_vpc" {
	name       = var.instance_name
	cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "CUSTOMIZE"
  vpc_restriction_effect = "DENY"
  vpc_id_set = [tencentcloud_vpc.pg_vpc2.id]
  ip_restriction_effect = "DENY"
  ip_set = ["192.168.0.0"]
}
```

Import

postgresql backup_download_restriction_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config backup_download_restriction_config_id
```