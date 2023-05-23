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

```hcl
resource "tencentcloud_postgresql_backup_download_restriction_config" "backup_download_restriction_config" {
  restriction_type       = ""
  vpc_restriction_effect = ""
  vpc_id_set             =
  ip_restriction_effect  = ""
  ip_set                 =
}
```

## Argument Reference

The following arguments are supported:

* `restriction_type` - (Required, String) Backup file download restriction type:NONE:Unlimited, both internal and external networks can be downloaded.INTRANET:Only intranet downloads are allowed.CUSTOMIZE:Customize the vpc or ip that limits downloads.
* `ip_restriction_effect` - (Optional, String) ip limit Strategy:ALLOW,DENY.
* `ip_set` - (Optional, Set: [`String`]) The list of ip&amp;#39;s that are allowed or denied to download backup files.
* `vpc_id_set` - (Optional, Set: [`String`]) The list of vpcIds that allow or deny downloading of backup files.
* `vpc_restriction_effect` - (Optional, String) vpc limit Strategy:ALLOW,DENY.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql backup_download_restriction_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_download_restriction_config.backup_download_restriction_config backup_download_restriction_config_id
```

