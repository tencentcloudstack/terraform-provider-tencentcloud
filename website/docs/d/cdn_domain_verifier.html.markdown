---
subcategory: "Content Delivery Network(CDN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdn_domain_verifier"
sidebar_current: "docs-tencentcloud-datasource-cdn_domain_verifier"
description: |-
  Provides a resource to check or create a cdn Domain Verify Record
---

# tencentcloud_cdn_domain_verifier

Provides a resource to check or create a cdn Domain Verify Record

~> **NOTE:**

## Example Usage

```hcl
data "tencentcloud_cdn_domain_verifier" "vr" {
  domain        = "www.examplexxx123.com"
  auto_verify   = true # auto create record if not verified
  freeze_record = true # once been freeze and verified, it will never be changed again
}

locals {
  recordValue = data.tencentcloud_cdn_domain_verifier.record
  recordType  = data.tencentcloud_cdn_domain_verifier.record_type
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Specify domain name, e.g. `www.examplexxx123.com`.
* `auto_verify` - (Optional, Bool) Specify whether to keep first create result instead of re-create again.
* `failed_reason` - (Optional, String) Indicates failed reason of verification.
* `freeze_record` - (Optional, Bool) Specify whether the verification record needs to be freeze instead of refresh every 8 hours, this used for domain verification.
* `result_output_file` - (Optional, String) Used for save result json.
* `verify_type` - (Optional, String) Specify verify type, values: `dns` (default), `file`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `file_verify_domains` - List of file verified domains.
* `file_verify_name` - Name of file verifications.
* `file_verify_url` - File verify URL guidance.
* `record_type` - Type of resolution.
* `record` - Resolution record value.
* `sub_domain` - Sub-domain resolution.
* `verify_result` - Verify result.


