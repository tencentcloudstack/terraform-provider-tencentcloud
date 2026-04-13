---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_remediation"
sidebar_current: "docs-tencentcloud-resource-config_remediation"
description: |-
  Provides a resource to create a Config remediation setting.
---

# tencentcloud_config_remediation

Provides a resource to create a Config remediation setting.

## Example Usage

```hcl
resource "tencentcloud_config_remediation" "example" {
  rule_id                 = "cr-Gxt8pzxgCVZJ0C95H1HO"
  remediation_type        = "SCF"
  remediation_template_id = "qcs::scf:ap-guangzhou:uin/100000005287:namespace/test/functions/my-remediation-func"
  invoke_type             = "MANUAL_EXECUTION"
  source_type             = "CUSTOM"
}
```

## Argument Reference

The following arguments are supported:

* `invoke_type` - (Required, String) Remediation execution mode. Valid values: MANUAL_EXECUTION (manual), AUTO_EXECUTION (automatic), NON_EXECUTION (disabled), NOT_CONFIG (not configured).
* `remediation_template_id` - (Required, String) Remediation template ID (e.g. SCF function resource path: qcs::scf:ap-guangzhou:uin/functions/xxx).
* `remediation_type` - (Required, String) Remediation type. Valid value: SCF (cloud function, custom remediation).
* `rule_id` - (Required, String, ForceNew) Config rule ID to bind the remediation setting to.
* `source_type` - (Optional, String) Template source. Valid value: CUSTOM (custom template).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `owner_uin` - Owner account UIN.
* `remediation_id` - Remediation setting ID.
* `remediation_source_type` - Remediation source type returned from API.


## Import

Config remediation can be imported using the id, e.g.

```
terraform import tencentcloud_config_remediation.example crr-lKj43O4nbSD78XYlvGS9
```

