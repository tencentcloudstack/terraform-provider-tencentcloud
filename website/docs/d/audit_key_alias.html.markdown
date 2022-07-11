---
subcategory: "Audit"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_audit_key_alias"
sidebar_current: "docs-tencentcloud-datasource-audit_key_alias"
description: |-
  Use this data source to query the key alias list specified with region supported by the audit.
---

# tencentcloud_audit_key_alias

Use this data source to query the key alias list specified with region supported by the audit.

## Example Usage

```hcl
data "tencentcloud_audit_key_alias" "all" {
  region = "ap-hongkong"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) Region.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `audit_key_alias_list` - List of available key alias supported by audit.
  * `key_alias` - Key alias.
  * `key_id` - Key ID.


