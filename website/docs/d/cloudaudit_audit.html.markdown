---
subcategory: "Cloud Audit(CloudAudit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cloudaudit_audit"
sidebar_current: "docs-tencentcloud-data_source-cloudaudit_audit"
description: |-
  Use this data source to query detailed information of a CloudAudit audit track set.
---

# tencentcloud_cloudaudit_audit

Use this data source to query detailed information of a CloudAudit audit track set.

## Example Usage

```hcl
data "tencentcloud_cloudaudit_audit" "audit" {
  audit_name = "tf-audit"
}
```

## Argument Reference

The following arguments are supported:

* `audit_name` - (Required, String) Name of the audit track set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audit_status` - Status of the audit track set. `1`: enabled, `0`: disabled.
* `cmq_queue_name` - Name of the CMQ queue.
* `cmq_region` - Region of the CMQ queue.
* `cos_bucket_name` - Name of the COS bucket.
* `cos_region` - Region of the COS bucket.
* `is_enable_cmq_notify` - Whether CMQ message notification is enabled. `1`: yes, `0`: no.
* `is_enable_kms_encry` - Whether KMS encryption is enabled. `1`: yes, `0`: no. When KMS encryption is enabled, data is encrypted before being delivered to COS.
* `key_id` - Global unique identifier of the CMK.
* `kms_alias` - Alias of the CMK.
* `kms_region` - Region of KMS.
* `log_file_prefix` - Log file prefix.
* `read_write_attribute` - Management event read/write attribute. `1`: read only, `2`: write only, `3`: all.


