---
subcategory: "Global Accelerator(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_accelerate_area"
sidebar_current: "docs-tencentcloud-resource-ga2_accelerate_area"
description: |-
  Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) accelerate area.
---

# tencentcloud_ga2_accelerate_area

Provides a resource to create a Tencent Cloud Global Accelerator V2 (GA2) accelerate area.

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

resource "tencentcloud_ga2_accelerate_area" "example" {
  global_accelerator_id = tencentcloud_ga2_global_accelerator.example.id
  accelerate_region     = "ap-guangzhou"
  bandwidth             = 10
  isp_type              = "BGP"
  ip_version            = "IPv4"
}
```

## Argument Reference

The following arguments are supported:

* `accelerate_region` - (Required, String, ForceNew) Acceleration region. Serves as the natural key used to resolve the acceleration region ID after creation. Cannot be modified after creation; modifying it forces a new resource.
* `global_accelerator_id` - (Required, String, ForceNew) Global accelerator instance ID this acceleration region belongs to.
* `bandwidth` - (Optional, Int) Acceleration bandwidth in Mbps.
* `ip_address` - (Optional, Set: [`String`]) Bound IP address list. Treated as an unordered set; HCL element order has no semantic meaning.
* `ip_version` - (Optional, String) IP version. Only `IPv4` is supported. Default: `IPv4`.
* `isp_type` - (Optional, String) ISP type. Valid values: `BGP` (BGP), `STATIC_IP` (multi-ISP static IP), `QUALITY_BGP` (premium BGP). Default: `BGP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `accelerator_area_id` - Acceleration region ID.
* `ip_address_info_set` - IP address information list.
  * `ip_address` - IP address.
  * `isp_type` - ISP type of the IP address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `update` - (Defaults to `5m`) Used when updating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

GA2 accelerate area can be imported using the composite id `<global_accelerator_id>#<accelerator_area_id>`, e.g.

```
terraform import tencentcloud_ga2_accelerate_area.example ga-jg9gepn0#area-jrsub43y
```

