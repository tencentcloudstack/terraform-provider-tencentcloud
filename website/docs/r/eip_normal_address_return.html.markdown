---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip_normal_address_return"
sidebar_current: "docs-tencentcloud-resource-eip_normal_address_return"
description: |-
  Provides a resource to create a vpc normal_address_return
---

# tencentcloud_eip_normal_address_return

Provides a resource to create a vpc normal_address_return

## Example Usage

```hcl
resource "tencentcloud_eip_normal_address_return" "normal_address_return" {
  address_ips =
}
```

## Argument Reference

The following arguments are supported:

* `address_ips` - (Optional, Set: [`String`], ForceNew) The IP address of the EIP, example: 101.35.139.183.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



