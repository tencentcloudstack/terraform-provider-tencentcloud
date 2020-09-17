---
subcategory: "Audit"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit"
sidebar_current: "docs-tencentcloud-resource-audit"
description: |-
  Provides a resource to create an audit.
---

# tencentcloud_audit

Provides a resource to create an audit.

## Example Usage

```hcl
resource "tencentcloud_audit" "foo" {
  name                 = "audittest"
  cos_bucket           = "test"
  cos_region           = "ap-hongkong"
  log_file_prefix      = "test"
  audit_switch         = true
  read_write_attribute = 3
}
```

## Argument Reference

The following arguments are supported:

* `audit_switch` - (Required) Indicate whether to turn on audit logging or not.
* `cos_bucket` - (Required) Name of the cos bucket to save audit log. Caution: the validation of existing cos bucket will not be checked by terraform.
* `cos_region` - (Required) Region of the cos bucket.
* `name` - (Required, ForceNew) Name of audit. Valid length ranges from 3 to 128. Only alpha character or numbers or `_` supported.
* `read_write_attribute` - (Required) Event attribute filter. 1 for readonly, 2 for writeonly, 3 for all.
* `enable_kms_encry` - (Optional) Indicate whether the log is encrypt with KMS algorithm or not.
* `key_id` - (Optional) Existing CMK unique key. This field can be get by data source `tencentcloud_audit_key_alias`. Caution: the region of the KMS must be as same as the `cos_region`.
* `log_file_prefix` - (Optional) The log file name prefix. The length ranges from 3 to 40. If not set, the account ID will be the log file prefix.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Audit can be imported using the id, e.g.

```
$ terraform import tencentcloud_audit.foo audit-test
```

