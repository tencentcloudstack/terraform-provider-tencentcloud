---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elastic_public_ipv6_attachment"
sidebar_current: "docs-tencentcloud-resource-elastic_public_ipv6_attachment"
description: |-
  Provides a resource to create a vpc elastic_public_ipv6_attachment
---

# tencentcloud_elastic_public_ipv6_attachment

Provides a resource to create a vpc elastic_public_ipv6_attachment

## Example Usage

```hcl
resource "tencentcloud_elastic_public_ipv6_attachment" "elastic_public_ipv6_attachment" {
  ipv6_address_id      = "eipv6-xxxxxx"
  network_interface_id = "eni-xxxxxx"
  private_ipv6_address = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `ipv6_address_id` - (Required, String, ForceNew) Elastic IPv6 unique ID, EIPv6 unique ID is like eipv6-11112222.
* `keep_bind_with_eni` - (Optional, Bool) Whether to keep the Elastic Network Interface bound when unbinding.
* `network_interface_id` - (Optional, String, ForceNew) Elastic Network Interface ID to bind. Elastic Network Interface ID is like eni-11112222. NetworkInterfaceId and InstanceId cannot be specified simultaneously. The Elastic Network Interface ID can be queried by logging in to the console, or obtained through the networkInterfaceId in the return value of the DescribeNetworkInterfaces interface.
* `private_ipv6_address` - (Optional, String, ForceNew) The intranet IPv6 to bind. If NetworkInterfaceId is specified, PrivateIPv6Address must also be specified, which means that the EIP is bound to the specified private network IP of the specified Elastic Network Interface. Also ensure that the specified PrivateIPv6Address is an intranet IPv6 on the specified NetworkInterfaceId. The intranet IPv6 of the specified Elastic Network Interface can be queried by logging in to the console, or obtained through the Ipv6AddressSet.Address in the return value of the DescribeNetworkInterfaces interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc elastic_public_ipv6_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_elastic_public_ipv6_attachment.elastic_public_ipv6_attachment elastic_public_ipv6_attachment_id
```

