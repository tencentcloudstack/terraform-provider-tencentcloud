---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_classic_elastic_public_ipv6s"
sidebar_current: "docs-tencentcloud-datasource-classic_elastic_public_ipv6s"
description: |-
  Use this data source to query detailed information of vpc classic_elastic_public_ipv6s
---

# tencentcloud_classic_elastic_public_ipv6s

Use this data source to query detailed information of vpc classic_elastic_public_ipv6s

## Example Usage

```hcl
data "tencentcloud_classic_elastic_public_ipv6s" "classic_elastic_public_ipv6s" {
  ip6_address_ids = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) The upper limit for `Filters per request is 10, and the upper limit for`Filter.Values` is 100. Parameters do not support specifying both AddressIds and Filters. The detailed filtering conditions are as follows:
  - address-ip: filter according to IPV6 IP address.
  - network-interface-id: filter according to the unique ID of the Elastic Network Interface.
* `ip6_address_ids` - (Optional, Set: [`String`]) List of unique IDs that identify IPV6. The IPV6 unique ID is shaped like `eip-11112222`. Parameters do not support specifying both `Ip6AddressIds` and `Filters`.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Property name. If there are multiple Filters, the relationship between Filters is a logical AND relationship.
* `values` - (Required, Set) Attribute value. If there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR relationship. When the value type is a Boolean type, the value can be directly taken to the string TRUE or FALSE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address_set` - List of IPV6 details.


