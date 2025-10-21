---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_file_download_url"
sidebar_current: "docs-tencentcloud-datasource-dcdb_file_download_url"
description: |-
  Use this data source to query detailed information of dcdb file_download_url
---

# tencentcloud_dcdb_file_download_url

Use this data source to query detailed information of dcdb file_download_url

## Example Usage

```hcl
data "tencentcloud_dcdb_file_download_url" "file_download_url" {
  instance_id = local.dcdb_id
  shard_id    = "shard-1b5r04az"
  file_path   = "/cos_backup/test.txt"
}
```

## Argument Reference

The following arguments are supported:

* `file_path` - (Required, String) Unsigned file path.
* `instance_id` - (Required, String) Instance ID.
* `shard_id` - (Required, String) Instance Shard ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `pre_signed_url` - Signed download URL.


