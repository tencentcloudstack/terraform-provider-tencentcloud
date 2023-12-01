---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_backup_download_restriction_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_backup_download_restriction_config"
description: |-
  Provides a resource to create a postgresql backup_download_restriction_config
---

# tencentcloud_postgresql_backup_download_restriction_config

Provides a resource to create a postgresql backup_download_restriction_config

## Example Usage

### Unlimit the restriction of the backup file download.

```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "NONE"
}
```

### Set the download only to allow the intranet downloads.

```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type = "INTRANET"
}
```

### Restrict the backup file download by customizing.

```hcl
resource "tencentcloud_vpc" "pg_vpc" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type       = "CUSTOMIZE"
  vpc_restriction_effect = "DENY"
  vpc_id_set             = [tencentcloud_vpc.pg_vpc2.id]
  ip_restriction_effect  = "DENY"
  ip_set                 = ["192.168.0.0"]
}
```

## Argument Reference

The following arguments are supported:

* `restriction_type` - (Required, String) Backup file download restriction type: NONE:Unlimited, both internal and external networks can be downloaded. INTRANET:Only intranet downloads are allowed. CUSTOMIZE:Customize the vpc or ip that limits downloads.
* `ip_restriction_effect` - (Optional, String) ip limit Strategy: ALLOW, DENY.
* `ip_set` - (Optional, Set: [`String`]) The list of ips that are allowed or denied to download backup files.
* `vpc_id_set` - (Optional, Set: [`String`]) The list of vpcIds that allow or deny downloading of backup files.
* `vpc_restriction_effect` - (Optional, String) vpc limit Strategy: ALLOW, DENY.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql backup_download_restriction_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config backup_download_restriction_config_id
```

