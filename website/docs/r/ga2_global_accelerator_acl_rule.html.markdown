---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_global_accelerator_acl_rule"
sidebar_current: "docs-tencentcloud-resource-ga2_global_accelerator_acl_rule"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule.
---

# tencentcloud_ga2_global_accelerator_acl_rule

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule.

## Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"
}

resource "tencentcloud_ga2_global_accelerator_acl_rule" "example" {
  global_accelerator_id            = tencentcloud_ga2_global_accelerator.example.id
  global_accelerator_acl_policy_id = "aclpol-xxxxxxxx"
  protocol                         = "TCP"
  port                             = "80"
  source_cidr_block                = "10.0.0.0/24"
  policy                           = "ACCEPT"
  description                      = "tf example acl rule"
}
```

## Argument Reference

The following arguments are supported:

* `global_accelerator_acl_policy_id` - (Required, String, ForceNew) ACL policy ID.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.
* `policy` - (Required, String) Action. Valid values: `ACCEPT` (allow), `DROP` (deny).
* `port` - (Required, String) Port.
* `protocol` - (Required, String) Protocol. Valid values: `TCP`, `UDP`, `ALL`.
* `source_cidr_block` - (Required, String) Source CIDR block.
* `description` - (Optional, String) Description. Maximum length is 100 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `global_accelerator_acl_rule_id` - ACL rule ID.
* `task_id` - Async task ID for the last operation on this resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 ACL rule can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>#<global_accelerator_acl_rule_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_rule.example ga-xxxxxxxx#aclpol-xxxxxxxx#aclrule-xxxxxxxx
```

