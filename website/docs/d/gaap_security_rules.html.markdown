---
subcategory: "Global Application Acceleration(GAAP)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_security_rules"
sidebar_current: "docs-tencentcloud-datasource-gaap_security_rules"
description: |-
  Use this data source to query security policy rule.
---

# tencentcloud_gaap_security_rules

Use this data source to query security policy rule.

## Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "ACCEPT"
}

resource "tencentcloud_gaap_security_rule" "foo" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  name      = "ci-test-gaap-s-rule"
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
  port      = "80"
}

data "tencentcloud_gaap_security_rules" "protocol" {
  policy_id = tencentcloud_gaap_security_policy.foo.id
  protocol  = tencentcloud_gaap_security_rule.foo.protocol
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) ID of the security policy to be queried.
* `action` - (Optional, String) Policy of the rule to be queried.
* `cidr_ip` - (Optional, String) A network address block of the request source to be queried.
* `name` - (Optional, String) Name of the security policy rule to be queried.
* `port` - (Optional, String) Port of the security policy rule to be queried.
* `protocol` - (Optional, String) Protocol of the security policy rule to be queried.
* `result_output_file` - (Optional, String) Used to save results.
* `rule_id` - (Optional, String) ID of the security policy rules to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rules` - An information list of security policy rule. Each element contains the following attributes:
  * `action` - Policy of the rule.
  * `cidr_ip` - A network address block of the request source.
  * `id` - ID of the security policy rule.
  * `name` - Name of the security policy rule.
  * `port` - Port of the security policy rule.
  * `protocol` - Protocol of the security policy rule.


