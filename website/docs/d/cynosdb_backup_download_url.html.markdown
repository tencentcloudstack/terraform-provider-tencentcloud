---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_backup_download_url"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_backup_download_url"
description: |-
  Use this data source to query detailed information of cynosdb backup_download_url
---

# tencentcloud_cynosdb_backup_download_url

Use this data source to query detailed information of cynosdb backup_download_url

## Example Usage

```hcl
data "tencentcloud_cynosdb_backup_download_url" "backup_download_url" {
  cluster_id = "cynosdbmysql-bws8h88b"
  backup_id  = 480782
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Required, Int) Backup ID.
* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `download_url` - Backup download address.


