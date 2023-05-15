---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_download_restriction"
sidebar_current: "docs-tencentcloud-resource-redis_backup_download_restriction"
description: |-
  Provides a resource to create a redis backup_download_restriction
---

# tencentcloud_redis_backup_download_restriction

Provides a resource to create a redis backup_download_restriction

## Example Usage

```hcl
resource "tencentcloud_redis_backup_download_restriction" "backup_download_restriction" {
  limit_type            = "Customize"
  vpc_comparison_symbol = "In"
  ip_comparison_symbol  = "In"
  limit_vpc {
    region   = "ap-guangzhou"
    vpc_list = [var.vpc_id]
  }
  limit_ip = ["10.1.1.12", "10.1.1.13"]
}
```

## Argument Reference

The following arguments are supported:

* `limit_type` - (Required, String) Types of network restrictions for downloading backup files:- NoLimit: There is no limit, and backup files can be downloaded from both Tencent Cloud and internal and external networks.- LimitOnlyIntranet: Only intranet addresses automatically assigned by Tencent Cloud can download backup files.- Customize: refers to a user-defined private network downloadable backup file.
* `ip_comparison_symbol` - (Optional, String) Identifies whether the customized LimitIP address can download the backup file.- In: Custom IP addresses are available for download.- NotIn: Custom IPs are not available for download.
* `limit_ip` - (Optional, Set: [`String`]) A custom VPC IP address for downloadable backup files.If the parameter LimitType is **Customize**, you need to configure this parameter.
* `limit_vpc` - (Optional, List) A custom VPC ID for a downloadable backup file.If the parameter LimitType is **Customize**, you need to configure this parameter.
* `vpc_comparison_symbol` - (Optional, String) This parameter only supports entering In, which means that the custom LimitVpc can download the backup file.

The `limit_vpc` object supports the following:

* `region` - (Required, String) Customize the region of the VPC to which the backup file is downloaded.
* `vpc_list` - (Required, Set) Customize the list of VPCs to download backup files.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

redis backup_download_restriction can be imported using the id, e.g.

```
terraform import tencentcloud_redis_backup_download_restriction.backup_download_restriction backup_download_restriction_id
```

