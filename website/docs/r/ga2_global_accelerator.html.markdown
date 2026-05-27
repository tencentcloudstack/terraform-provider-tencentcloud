---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_global_accelerator"
sidebar_current: "docs-tencentcloud-resource-ga2_global_accelerator"
description: |-
  Provides a resource to create a GA2 (Global Accelerator 2) global accelerator instance.
---

# tencentcloud_ga2_global_accelerator

Provides a resource to create a GA2 (Global Accelerator 2) global accelerator instance.

## Example Usage

```hcl
resource "tencentcloud_ga2_global_accelerator" "example" {
  name                      = "tf-example"
  instance_charge_type      = "POSTPAID"
  description               = "terraform example global accelerator"
  cross_border_type         = "HighQuality"
  cross_border_promise_flag = true

  tags = {
    createdBy = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cross_border_promise_flag` - (Optional, Bool) Flag indicating acceptance of cross-border service agreement. Must be set to `true` when using cross-border service.
* `cross_border_type` - (Optional, String) Cross-border type. Valid values: `HighQuality` (premium BGP-IP cross-border), `Unicom` (China Unicom dedicated line cross-border).
* `description` - (Optional, String) Description of the global accelerator instance. Maximum length is 100 bytes.
* `instance_charge_type` - (Optional, String, ForceNew) Billing mode. Valid values: `PREPAID` (prepaid, monthly subscription), `POSTPAID` (postpaid, pay-as-you-go). Default: `POSTPAID`. Currently only postpaid is supported.
* `name` - (Optional, String) Name of the global accelerator instance. Maximum length is 60 bytes.
* `tags` - (Optional, Map, ForceNew) Tag information.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname` - CNAME domain of the global accelerator instance.
* `create_time` - Creation time of the global accelerator instance.
* `ddos_id` - DDoS protection ID of the global accelerator instance.
* `state` - State of the global accelerator instance.
* `status` - Status of the global accelerator instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

GA2 global accelerator can be imported using the ID, e.g.

```
terraform import tencentcloud_ga2_global_accelerator.example ga2-xxxxxxxx
```

