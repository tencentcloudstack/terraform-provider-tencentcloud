---
subcategory: "Global Application Acceleration(GAAP)"
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

* `action` - (Required, String) Policy of the rule. Valid value: `ACCEPT` and `DROP`.
* `cidr_ip` - (Required, String) A network address block of the request source.
* `policy_id` - (Required, String, ForceNew) ID of the security policy.
* `name` - (Optional, String) Name of the security policy rule. Maximum length is 30.
* `port` - (Optional, String) Target port. Default value is `ALL`. Valid examples: `80`, `80,443` and `3306-20000`.
* `protocol` - (Optional, String) Protocol of the security policy rule. Default value is `ALL`. Valid value: `TCP`, `UDP` and `ALL`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSecurityRules](https://www.tencentcloud.com/document/product/608/33107) - This API is used to add security policy rules.
* Read: [DescribeSecurityRules](https://www.tencentcloud.com/document/product/608/33138) - This API is used to query the list of security rules based on security rule ID. It supports querying 1 to 20 security rules at a time.
* Update: [ModifySecurityRule](https://www.tencentcloud.com/document/product/608/33091) - This API is used to modify security policy rule names.
* Delete: [DeleteSecurityRules](https://www.tencentcloud.com/document/product/608/33105) - This API is used to delete security policy rules.

## Import

GAAP security rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_rule.foo sr-xxxxxxxx
```

