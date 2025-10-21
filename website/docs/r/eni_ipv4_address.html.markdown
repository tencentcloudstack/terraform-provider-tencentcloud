---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_ipv4_address"
sidebar_current: "docs-tencentcloud-resource-eni_ipv4_address"
description: |-
  Provides a resource to create a vpc eni ipv4 address
---

# tencentcloud_eni_ipv4_address

Provides a resource to create a vpc eni ipv4 address

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "example" {
  name        = "tf-example"
  description = "sg desc."
  project_id  = 0

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
  security_groups = [
    tencentcloud_security_group.example.id,
  ]
}

resource "tencentcloud_eni_ipv4_address" "example" {
  network_interface_id               = tencentcloud_eni.example.id
  qos_level                          = "DEFAULT"
  secondary_private_ip_address_count = 3
}
```

### Or

```hcl
resource "tencentcloud_eni_ipv4_address" "example" {
  network_interface_id = tencentcloud_eni.example.id
  private_ip_addresses {
    is_wan_ip_blocked  = false
    private_ip_address = "10.0.0.15"
    qos_level          = "DEFAULT"
  }

  private_ip_addresses {
    is_wan_ip_blocked  = false
    private_ip_address = "10.0.0.4"
    qos_level          = "DEFAULT"
  }
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

vpc eni ipv4 address can be imported using the id, e.g.

```
terraform import tencentcloud_eni_ipv4_address.example eni-65369ozn
```

