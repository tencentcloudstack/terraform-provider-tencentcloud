---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_bandwidth_package_resources"
sidebar_current: "docs-tencentcloud-resource-vpc_bandwidth_package_resources"
description: |-
  Provides a resource to create a vpc bandwidth_package_resources
---

# tencentcloud_vpc_bandwidth_package_resources

Provides a resource to create a vpc bandwidth_package_resources

## Example Usage

```hcl
resource "tencentcloud_vpc_bandwidth_package_resources" "bandwidth_package_resources" {
  resource_ids         = "lb-dv1ai6ma"
  bandwidth_package_id = "bwp-atmf0p9g"
  network_type         = "BGP"
  resource_type        = "LoadBalance"
  protocol             = ""
}
```

## Argument Reference

The following arguments are supported:

* `resource_ids` - (Required, String, ForceNew) The unique ID of the resource, currently supports EIP resources and LB resources, such as `eip-xxxx`, `lb-xxxx`.
* `bandwidth_package_id` - (Optional, String, ForceNew) Bandwidth package unique ID, in the form of `bwp-xxxx`.
* `network_type` - (Optional, String, ForceNew) Bandwidth packet type, currently supports `BGP` type, indicating that the internal resource is BGP IP.
* `protocol` - (Optional, String, ForceNew) Bandwidth packet protocol type. Currently `ipv4` and `ipv6` protocol types are supported.
* `resource_type` - (Optional, String, ForceNew) Resource types, including `Address`, `LoadBalance`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



