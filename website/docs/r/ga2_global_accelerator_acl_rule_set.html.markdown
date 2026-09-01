---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_global_accelerator_acl_rule_set"
sidebar_current: "docs-tencentcloud-resource-ga2_global_accelerator_acl_rule_set"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule set that manages the full collection of ACL rules under one ACL policy.
---

# tencentcloud_ga2_global_accelerator_acl_rule_set

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) ACL rule set that manages the full collection of ACL rules under one ACL policy.

~> **NOTE:** This resource must exclusive in one acl policy, do not declare additional acl rule resources of this acl policy elsewhere.

~> **NOTE:** The field length for `acl_entries` is subject to account resource quotas.

## Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                 = "tf-example"
  instance_charge_type = "POSTPAID"
  description          = "tf example global accelerator"

  tags = {
    createdBy = "Terraform"
  }
}

resource "tencentcloud_ga2_global_accelerator_acl_policy" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  default_action        = "ACCEPT"
  status                = "OPEN"
}

resource "tencentcloud_ga2_global_accelerator_acl_rule_set" "example" {
  global_accelerator_id            = tencentcloud_ga2_global_accelerator.example.id
  global_accelerator_acl_policy_id = tencentcloud_ga2_global_accelerator_acl_policy.example.global_accelerator_acl_policy_id

  acl_entries {
    protocol          = "TCP"
    port              = "80"
    source_cidr_block = "10.0.0.0/24"
    policy            = "ACCEPT"
    description       = "tf example acl rule 1"
  }

  acl_entries {
    protocol          = "UDP"
    port              = "443"
    source_cidr_block = "10.0.1.0/24"
    policy            = "DROP"
    description       = "tf example acl rule 2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Required, Set) The desired full set of ACL rules under the policy. Treated as an unordered set; HCL element order has no semantic meaning.
* `global_accelerator_acl_policy_id` - (Required, String, ForceNew) ACL policy ID that owns the rule set.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID.

The `acl_entries` object supports the following:

* `policy` - (Required, String) Action. Valid values: `ACCEPT` (allow), `DROP` (deny).
* `port` - (Required, String) Port.
* `protocol` - (Required, String) Protocol. Valid values: `TCP`, `UDP`.
* `source_cidr_block` - (Required, String) Source CIDR block.
* `description` - (Optional, String) Description. Maximum length is 100 bytes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.

The `acl_entries` object exports the following:

* `global_accelerator_acl_rule_id` - ACL rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `10m`) Used when deleting the resource.

## Import

GA2 ACL rule set can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_rule_set.example ga-cmbzp36q#sp-jz94sb2t
```

