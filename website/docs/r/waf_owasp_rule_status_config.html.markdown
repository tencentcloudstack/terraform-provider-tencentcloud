---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_owasp_rule_status_config"
sidebar_current: "docs-tencentcloud-resource-waf_owasp_rule_status_config"
description: |-
  Provides a resource to create a WAF owasp rule status config
---

# tencentcloud_waf_owasp_rule_status_config

Provides a resource to create a WAF owasp rule status config

## Example Usage

```hcl
resource "tencentcloud_waf_owasp_rule_status_config" "example" {
  domain      = "demo.com"
  rule_id     = "106251141"
  rule_status = 1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain name.
* `rule_id` - (Required, String, ForceNew) Rule ID.
* `rule_status` - (Required, Int) Rule switch. valid values: 0 (disabled), 1 (enabled), 2 (observation only).
* `reason` - (Optional, Int) Reason for modification. valid values: 0: none (compatibility record is empty). 1: avoid false positives due to business characteristics. 2: reporting of rule-based false positives. 3: gray release of core business rules. 4: others.
* `type_id` - (Optional, Int) If reverse requires the input of data type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cve_id` - CVE ID.
* `description` - Rule description.
* `level` - Protection level of the rule. valid values: 100 (loose), 200 (normal), 300 (strict), 400 (ultra-strict).
* `locked` - Whether the user is locked.
* `vul_level` - Threat level. valid values: 0 (unknown), 100 (low risk), 200 (medium risk), 300 (high risk), 400 (critical).


## Import

WAF owasp rule status config can be imported using the domain#ruleId, e.g.

```
terraform import tencentcloud_waf_owasp_rule_status_config.example demo.com#106251141
```

