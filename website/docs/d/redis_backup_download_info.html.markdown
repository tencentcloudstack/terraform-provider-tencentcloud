---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_backup_download_info"
sidebar_current: "docs-tencentcloud-datasource-redis_backup_download_info"
description: |-
  Use this data source to query detailed information of redis backup_download_info
---

# tencentcloud_redis_backup_download_info

Use this data source to query detailed information of redis backup_download_info

## Example Usage

```hcl
data "tencentcloud_redis_backup_download_info" "backup_download_info" {
  instance_id = "crs-iw7d9wdd"
  backup_id   = "641186639-8362913-1516672770"
  # limit_type = "NoLimit"
  # vpc_comparison_symbol = "In"
  # ip_comparison_symbol = "In"
  # limit_vpc {
  # 	region = "ap-guangzhou"
  # 	vpc_list = [""]
  # }
  # limit_ip = [""]
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required, String) The backup ID, which can be accessed via [DescribeInstanceBackups](https://cloud.tencent.com/document/product/239/20011) interface returns the parameter RedisBackupSet to get.
* `instance_id` - (Required, String) The ID of instance.
* `ip_comparison_symbol` - (Optional, String) Identifies whether the customized LimitIP address can download the backup file.- In: Custom IP addresses are available for download.- NotIn: Custom IPs are not available for download.
* `limit_ip` - (Optional, Set: [`String`]) A custom VPC IP address for downloadable backup files.If the parameter LimitType is **Customize**, you need to configure this parameter.
* `limit_type` - (Optional, String) Types of network restrictions for downloading backup files:- NoLimit: There is no limit, and backup files can be downloaded from both Tencent Cloud and internal and external networks.- LimitOnlyIntranet: Only intranet addresses automatically assigned by Tencent Cloud can download backup files.- Customize: refers to a user-defined private network downloadable backup file.
* `limit_vpc` - (Optional, List) A custom VPC ID for a downloadable backup file.If the parameter LimitType is **Customize**, you need to configure this parameter.
* `result_output_file` - (Optional, String) Used to save results.
* `vpc_comparison_symbol` - (Optional, String) This parameter only supports entering In, which means that the custom LimitVpc can download the backup file.

The `limit_vpc` object supports the following:

* `region` - (Required, String) Customize the region of the VPC to which the backup file is downloaded.
* `vpc_list` - (Required, Set) Customize the list of VPCs to download backup files.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_infos` - A list of backup file information.
  * `download_url` - Backup file download address on the Internet (6 hours).
  * `file_name` - Backup file name.
  * `file_size` - The backup file size is in unit B, if it is 0, it is invalid.
  * `inner_download_url` - Backup file intranet download address (6 hours).


