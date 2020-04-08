---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gaap_security_rule"
sidebar_current: "docs-tencentcloud-resource-gaap_security_rule"
description: |-
  Provides a resource to create a security policy rule.
---

# tencentcloud_gaap_security_rule

Provides a resource to create a security policy rule.

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
  cidr_ip   = "1.1.1.1"
  action    = "ACCEPT"
  protocol  = "TCP"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, ForceNew) Policy of the rule, the available values include `ACCEPT` and `DROP`.
* `cidr_ip` - (Required, ForceNew) A network address block of the request source.
* `policy_id` - (Required, ForceNew) ID of the security policy.
* `name` - (Optional) Name of the security policy rule. Maximum length is 30.
* `port` - (Optional, ForceNew) Target port. Default value is `ALL`, the available values include `80`, `80,443` and `3306-20000`.
* `protocol` - (Optional, ForceNew) Protocol of the security policy rule. Default value is `ALL`, the available values include `TCP`, `UDP` and `ALL`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

GAAP security rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_rule.foo sr-xxxxxxxx
```

