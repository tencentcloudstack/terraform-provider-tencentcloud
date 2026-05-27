---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_global_accelerator"
sidebar_current: "docs-tencentcloud-resource-ga2_global_accelerator"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) instance.
---

# tencentcloud_ga2_global_accelerator

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) instance.

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
```

### Cross-border global accelerator

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                      = "tf-example"
  instance_charge_type      = "POSTPAID"
  description               = "tf example cross-border accelerator"
  cross_border_type         = "HighQuality"
  cross_border_promise_flag = true

  tags = {
    createdBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cross_border_promise_flag` - (Optional, Bool) Whether the cross-border service commitment letter is signed. `true` indicates signed. Required when `cross_border_type` is set.
* `cross_border_type` - (Optional, String) Cross-border type. Valid values: `HighQuality` (premium BGP-IP cross-border), `Unicom` (Unicom dedicated line cross-border).
* `description` - (Optional, String) Global accelerator instance description. Maximum length is 100 bytes.
* `instance_charge_type` - (Optional, String, ForceNew) Billing mode. `PREPAID` for monthly subscription, `POSTPAID` for pay-as-you-go. Default: `POSTPAID`. Currently only `POSTPAID` is supported. Cannot be changed after creation; modifying this attribute forces a new resource.
* `name` - (Optional, String) Global accelerator instance name. Maximum length is 60 bytes.
* `tags` - (Optional, Map) Tag key-value pairs to attach to the instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `accelerator_area_counts` - Number of acceleration regions attached to this global accelerator instance.
* `cname` - Acceleration domain (CNAME) assigned to the instance.
* `create_time` - Creation time of the global accelerator instance.
* `ddos_id` - Associated anti-DDoS instance ID.
* `global_accelerator_id` - Global accelerator instance ID.
* `listener_counts` - Number of listeners attached to this global accelerator instance.
* `state` - Provisioning state of the global accelerator instance.
* `status` - Operational status of the global accelerator instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 global accelerator instance can be imported using the id, e.g.

```
terraform import tencentcloud_ga2_global_accelerator.example ga-ar31grog
```

