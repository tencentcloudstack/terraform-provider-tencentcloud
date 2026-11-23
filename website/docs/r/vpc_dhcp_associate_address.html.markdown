---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_dhcp_associate_address"
sidebar_current: "docs-tencentcloud-resource-vpc_dhcp_associate_address"
description: |-
  Provides a resource to create a vpc dhcp associate address
---

# tencentcloud_vpc_dhcp_associate_address

Provides a resource to create a vpc dhcp associate address

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_dhcp_ip" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
  dhcp_ip_name = "tf-example"
}

resource "tencentcloud_eip" "eip" {
  name = "example-eip"
}

resource "tencentcloud_vpc_dhcp_associate_address" "example" {
  dhcp_ip_id = tencentcloud_vpc_dhcp_ip.example.id
  address_ip = tencentcloud_eip.eip.public_ip
}
```

## Argument Reference

The following arguments are supported:

* `address_ip` - (Required, String, ForceNew) Elastic public network `IP`. Must be `EIP` not bound to `DhcpIp`.
* `dhcp_ip_id` - (Required, String, ForceNew) `DhcpIp` unique `ID`, like: `dhcpip-9o233uri`. Must be a `DhcpIp` that is not bound to `EIP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AssociateDhcpIpWithAddressIp](https://www.tencentcloud.com/search/?wd=AssociateDhcpIpWithAddressIp&type=document&lang=en)
* Read: [DescribeDhcpIps](https://www.tencentcloud.com/search/?wd=DescribeDhcpIps&type=document&lang=en)
* Delete: [DisassociateDhcpIpWithAddressIp](https://www.tencentcloud.com/search/?wd=DisassociateDhcpIpWithAddressIp&type=document&lang=en)

## Import

vpc dhcp_associate_address can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dhcp_associate_address.dhcp_associate_address dhcp_associate_address_id
```

