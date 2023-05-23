---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_backup_download_restriction"
sidebar_current: "docs-tencentcloud-resource-mysql_backup_download_restriction"
description: |-
  Provides a resource to create a mysql backup_download_restriction
---

# tencentcloud_mysql_backup_download_restriction

Provides a resource to create a mysql backup_download_restriction

## Example Usage

```hcl
resource "tencentcloud_mysql_backup_download_restriction" "backup_download_restriction" {
  limit_type            = "Customize"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol  = "In"
  limit_vpc {
    region   = "ap-guangzhou"
    vpc_list = ["vpc-4owdpnwr"]
  }
  limit_ip = ["127.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `limit_type` - (Required, String) NoLimit No limit, both internal and external networks can be downloaded; LimitOnlyIntranet Only intranet can be downloaded; Customize user-defined vpc:ip can be downloaded. LimitVpc and LimitIp can be set only when the value is Customize.
* `ip_comparison_symbol` - (Optional, String) In: The specified ip can be downloaded; NotIn: The specified ip cannot be downloaded. The default is In.
* `limit_ip` - (Optional, Set: [`String`]) ip settings to limit downloads.
* `limit_vpc` - (Optional, List) vpc settings to limit downloads.
* `vpc_comparison_symbol` - (Optional, String) This parameter only supports In, which means that the vpc specified by LimitVpc can be downloaded. The default is In.

The `limit_vpc` object supports the following:

* `region` - (Required, String) Restrict downloads from regions. Currently only the current region is supported.
* `vpc_list` - (Required, Set) List of vpcs to limit downloads.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql backup_download_restriction can be imported using the "BackupDownloadRestriction", as follows.

```
terraform import tencentcloud_mysql_backup_download_restriction.backup_download_restriction BackupDownloadRestriction
```

