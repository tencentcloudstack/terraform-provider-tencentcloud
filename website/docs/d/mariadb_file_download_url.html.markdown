---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_file_download_url"
sidebar_current: "docs-tencentcloud-datasource-mariadb_file_download_url"
description: |-
  Use this data source to query detailed information of mariadb file_download_url
---

# tencentcloud_mariadb_file_download_url

Use this data source to query detailed information of mariadb file_download_url

## Example Usage

```hcl
data "tencentcloud_mariadb_file_download_url" "file_download_url" {
  instance_id = "tdsql-9vqvls95"
  file_path   = "/cos_backup/test.txt"
}
```

## Argument Reference

The following arguments are supported:

* `file_path` - (Required, String) Unsigned file path.
* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `pre_signed_url` - Signed download URL.


