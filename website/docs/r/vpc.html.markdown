---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc"
sidebar_current: "docs-tencentcloud-resource-vpc"
description: |-
  Provide a resource to create a VPC.
---

# tencentcloud_vpc

Provide a resource to create a VPC.

## Example Usage

```hcl
resource "tencentcloud_vpc" "foo" {
    name = "ci-temp-test-updated"
    cidr_block = "10.0.0.0/16"
	dns_servers=["119.29.29.29","8.8.8.8"]
	is_multicast=false
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).
* `name` - (Required) The name of the VPC.
* `dns_servers` - (Optional) The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.
* `is_multicast` - (Optional) Indicates whether VPC multicast is enabled. The default value is 'true'.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Creation time of VPC.
* `is_default` - Indicates whether it is the default VPC for this region.


## Import

Vpc instance can be imported, e.g.

```hcl
$ terraform import tencentcloud_vpc.test vpc-id
```

