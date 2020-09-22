---
subcategory: "Audit"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audits"
sidebar_current: "docs-tencentcloud-datasource-audits"
description: |-
  Use this data source to query detailed information of audits.
---

# tencentcloud_audits

Use this data source to query detailed information of audits.

## Example Usage

```hcl
data "tencentcloud_audits" "audits" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the audits.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audit_list` - Information list of the dedicated audits.
  * `audit_switch` - Indicate whether audit start logging or not.
  * `cos_bucket` - Cos bucket name where audit save logs.
  * `id` - Id of the audit.
  * `log_file_prefix` - Prefix of the log file of the audit.
  * `name` - Name of the audit.


