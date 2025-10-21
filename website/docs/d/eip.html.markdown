---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip"
sidebar_current: "docs-tencentcloud-datasource-eip"
description: |-
  Provides an available EIP for the user.
---

# tencentcloud_eip

Provides an available EIP for the user.

The EIP data source fetch proper EIP from user's EIP pool.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_eips.

## Example Usage

```hcl
data "tencentcloud_eip" "my_eip" {
  filter {
    name   = "address-status"
    values = ["UNBIND"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, Set) One or more name/value pairs to filter.
* `include_arrears` - (Optional, Bool) Whether the IP is arrears.
* `include_blocked` - (Optional, Bool) Whether the IP is blocked.

The `filter` object supports the following:

* `name` - (Required, String) Key of the filter, valid keys: `address-id`,`address-name`,`address-ip`.
* `values` - (Required, List) Value of the filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - An EIP id indicate the uniqueness of a certain EIP,  which can be used for instance binding or network interface binding.
* `public_ip` - An public IP address for the EIP.
* `status` - The status of the EIP, there are several status like `BIND`, `UNBIND`, and `BIND_ENI`.


