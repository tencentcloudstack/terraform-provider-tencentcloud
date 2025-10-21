---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_binlog_download_url"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_binlog_download_url"
description: |-
  Use this data source to query detailed information of cynosdb binlog_download_url
---

# tencentcloud_cynosdb_binlog_download_url

Use this data source to query detailed information of cynosdb binlog_download_url

## Example Usage

```hcl
data "tencentcloud_cynosdb_binlog_download_url" "binlog_download_url" {
  cluster_id = "cynosdbmysql-bws8h88b"
  binlog_id  = 6202249
}
```

## Argument Reference

The following arguments are supported:

* `binlog_id` - (Required, Int) Binlog file ID.
* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `download_url` - Download address.


