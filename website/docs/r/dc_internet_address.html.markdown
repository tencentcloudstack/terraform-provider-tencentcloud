---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_internet_address"
sidebar_current: "docs-tencentcloud-resource-dc_internet_address"
description: |-
  Provides a resource to create a dc internet_address
---

# tencentcloud_dc_internet_address

Provides a resource to create a dc internet_address

## Example Usage

```hcl
resource "tencentcloud_dc_internet_address" "internet_address" {
  mask_len   = 30
  addr_type  = 2
  addr_proto = 0
}
```

## Argument Reference

The following arguments are supported:

* `addr_proto` - (Required, Int, ForceNew) 0: IPv4, 1: IPv6.
* `addr_type` - (Required, Int, ForceNew) 0: BGP, 1: china telecom, 2: china mobile, 3: china unicom.
* `mask_len` - (Required, Int, ForceNew) CIDR address mask.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dc internet_address can be imported using the id, e.g.

```
terraform import tencentcloud_dc_internet_address.internet_address internet_address_id
```

