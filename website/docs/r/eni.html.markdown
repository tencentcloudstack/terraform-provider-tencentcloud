---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni"
sidebar_current: "docs-tencentcloud-resource-eni"
description: |-
  Provides a resource to create an ENI.
---

# tencentcloud_eni

Provides a resource to create an ENI.

## Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  availability_zone = "ap-guangzhou-3"
  name              = "ci-test-eni-subnet"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "foo" {
  name        = "ci-test-eni"
  vpc_id      = tencentcloud_vpc.foo.id
  subnet_id   = tencentcloud_subnet.foo.id
  description = "eni desc"
  ipv4_count  = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the ENI, maximum length 60.
* `subnet_id` - (Required, ForceNew) ID of the subnet within this vpc.
* `vpc_id` - (Required, ForceNew) ID of the vpc.
* `description` - (Optional) Description of the ENI, maximum length 60.
* `ipv4_count` - (Optional) The number of intranet IPv4s. When it is greater than 1, there is only one primary intranet IP. The others are auxiliary intranet IPs, which conflict with `ipv4s`.
* `ipv4s` - (Optional) Applying for intranet IPv4s collection, conflict with `ipv4_count`. When there are multiple ipv4s, can only be one primary IP, and the maximum length of the array is 30. Each element contains the following attributes:
* `security_groups` - (Optional) A set of security group IDs.
* `tags` - (Optional) Tags of the ENI.

The `ipv4s` object supports the following:

* `ip` - (Required) Intranet IP.
* `primary` - (Required) Indicates whether the IP is primary.
* `description` - (Optional) Description of the IP, maximum length 25.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the ENI.
* `ipv4_info` - An information list of IPv4s. Each element contains the following attributes:
  * `description` - Description of the IP.
  * `ip` - Intranet IP.
  * `primary` - Indicates whether the IP is primary.
* `mac` - MAC address.
* `primary` - Indicates whether the IP is primary.
* `state` - State of the ENI.


## Import

ENI can be imported using the id, e.g.

```
  $ terraform import tencentcloud_eni.foo eni-qka182br
```

