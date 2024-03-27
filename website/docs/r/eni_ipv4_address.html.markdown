---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_ipv4_address"
sidebar_current: "docs-tencentcloud-resource-eni_ipv4_address"
description: |-
  Provides a resource to create a vpc eni_ipv4_address
---

# tencentcloud_eni_ipv4_address

Provides a resource to create a vpc eni_ipv4_address

## Example Usage

```hcl
data "tencentcloud_enis" "eni" {
  name = "Primary ENI"
}

resource "tencentcloud_eni_ipv4_address" "eni_ipv4_address" {
  network_interface_id               = data.tencentcloud_enis.eni.enis.0.id
  secondary_private_ip_address_count = 3
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required, String, ForceNew) The ID of the ENI instance, such as `eni-m6dyj72l`.
* `private_ip_addresses` - (Optional, Set, ForceNew) The information on private IP addresses, of which you can specify a maximum of 10 at a time. You should provide either this parameter or SecondaryPrivateIpAddressCount, or both.
* `qos_level` - (Optional, String, ForceNew) IP service level. It is used together with `SecondaryPrivateIpAddressCount`. Values: PT`(Gold), `AU`(Silver), `AG `(Bronze) and DEFAULT (Default).
* `secondary_private_ip_address_count` - (Optional, Int, ForceNew) The number of newly-applied private IP addresses. You should provide either this parameter or PrivateIpAddresses, or both. The total number of private IP addresses cannot exceed the quota.

The `private_ip_addresses` object supports the following:

* `private_ip_address` - (Required, String, ForceNew) Private IP address.
* `address_id` - (Optional, String, ForceNew) EIP instance ID, such as `eip-11112222`.
* `description` - (Optional, String, ForceNew) Private IP description.
* `is_wan_ip_blocked` - (Optional, Bool, ForceNew) Whether the public IP is blocked.
* `primary` - (Optional, Bool, ForceNew) Whether it is a primary IP.
* `public_ip_address` - (Optional, String, ForceNew) Public IP address.
* `qos_level` - (Optional, String, ForceNew) IP service level. Values: PT` (Gold), `AU` (Silver), `AG `(Bronze) and DEFAULT` (Default).
* `state` - (Optional, String, ForceNew) IP status: `PENDING`: Creating, `MIGRATING`: Migrating, `DELETING`: Deleting, `AVAILABLE`: Available.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc eni_ipv4_address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv4_address.eni_ipv4_address eni_id
```

