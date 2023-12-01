---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_security_audit_log_download_urls"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_security_audit_log_download_urls"
description: |-
  Use this data source to query detailed information of dbbrain security_audit_log_download_urls
---

# tencentcloud_dbbrain_security_audit_log_download_urls

Use this data source to query detailed information of dbbrain security_audit_log_download_urls

## Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "%s"
  start_time         = "%s"
  end_time           = "%s"
  product            = "mysql"
  danger_levels      = [0, 1, 2]
}

data "tencentcloud_dbbrain_security_audit_log_download_urls" "test" {
  sec_audit_group_id = "%s"
  async_request_id   = tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id
  product            = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `async_request_id` - (Required, Int) Asynchronous task ID.
* `product` - (Required, String) Service product type, supported values: `mysql` - ApsaraDB for MySQL.
* `sec_audit_group_id` - (Required, String) Security audit group Id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `urls` - List of COS links to export results. When the result set is large, it may be divided into multiple urls for download.


