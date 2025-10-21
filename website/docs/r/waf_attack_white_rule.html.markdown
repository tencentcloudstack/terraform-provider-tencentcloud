---
subcategory: "Web Application Firewall(WAF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_attack_white_rule"
sidebar_current: "docs-tencentcloud-resource-waf_attack_white_rule"
description: |-
  Provides a resource to create a WAF attack white rule
---

# tencentcloud_waf_attack_white_rule

Provides a resource to create a WAF attack white rule

## Example Usage

### Using type_ids

```hcl
resource "tencentcloud_waf_attack_white_rule" "example" {
  domain = "www.demo.com"
  name   = "tf-example"
  status = 1
  mode   = 0
  rules {
    match_field   = "IP"
    match_method  = "ipmatch"
    match_content = "1.1.1.1"
  }

  rules {
    match_field   = "Referer"
    match_method  = "eq"
    match_content = "referer content"
  }

  rules {
    match_field   = "URL"
    match_method  = "contains"
    match_content = "/prefix"
  }

  rules {
    match_field   = "HTTP_METHOD"
    match_method  = "neq"
    match_content = "POST"
  }

  rules {
    match_field   = "GET"
    match_method  = "ncontains"
    match_content = "value"
    match_params  = "key"
  }

  type_ids = [
    "010000000",
    "020000000",
    "030000000",
    "040000000",
    "050000000",
    "060000000",
    "090000000",
    "110000000"
  ]
}
```

### Using signature_ids

```hcl
resource "tencentcloud_waf_attack_white_rule" "example" {
  domain = "www.demo.com"
  name   = "tf-example"
  status = 0
  mode   = 1
  rules {
    match_field   = "IP"
    match_method  = "ipmatch"
    match_content = "1.1.1.1"
  }

  signature_ids = [
    "60270036",
    "10000047"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) Domain.
* `rules` - (Required, List) Rule list.
* `status` - (Required, Int) Rule status.
* `mode` - (Optional, Int) 0: Whiten according to a specific rule ID, 1: Whiten according to the rule type.
* `name` - (Optional, String) Rule name.
* `signature_ids` - (Optional, Set: [`String`]) Whitelist of rule IDs.
* `type_ids` - (Optional, Set: [`String`]) The whitened category rule ID.

The `rules` object supports the following:

* `match_content` - (Required, String) Matching content.
* `match_field` - (Required, String) Matching domains.
* `match_method` - (Required, String) Matching method.
* `match_params` - (Optional, String) Matching params.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `rule_id` - Rule ID.


## Import

WAF attack white rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_attack_white_rule.example www.demo.com#38562
```

