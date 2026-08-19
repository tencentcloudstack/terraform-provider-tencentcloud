---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_global_accelerator_acl_policy"
sidebar_current: "docs-tencentcloud-resource-ga2_global_accelerator_acl_policy"
description: |-
  Provides a resource to create a GA2 (Global Accelerator V2) global accelerator ACL policy
---

# tencentcloud_ga2_global_accelerator_acl_policy

Provides a resource to create a GA2 (Global Accelerator V2) global accelerator ACL policy

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
```

## Argument Reference

The following arguments are supported:

* `default_action` - (Required, String, ForceNew) Default traffic action. Enumerated values: `ACCEPT` (allow all traffic by default), `DROP` (deny all traffic by default).
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this ACL policy belongs to.
* `status` - (Optional, String) ACL policy state. Enumerated values: `OPEN` (enabled), `CLOSE` (disabled). Default is `CLOSE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `global_accelerator_acl_policy_id` - ACL policy ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 global accelerator ACL policy can be imported using the composite id `<global_accelerator_id>#<global_accelerator_acl_policy_id>`, e.g.

```
terraform import tencentcloud_ga2_global_accelerator_acl_policy.example ga-ntc1iaco#sp-nseaj82f
```

