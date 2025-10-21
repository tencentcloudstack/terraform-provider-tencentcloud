---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elastic_public_ipv6s"
sidebar_current: "docs-tencentcloud-datasource-elastic_public_ipv6s"
description: |-
  Use this data source to query detailed information of vpc elastic_public_ipv6s
---

# tencentcloud_elastic_public_ipv6s

Use this data source to query detailed information of vpc elastic_public_ipv6s

## Example Usage

```hcl
data "tencentcloud_elastic_public_ipv6s" "elastic_public_ipv6s" {
  ipv6_address_ids = ["xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) The detailed filter conditions are as follows:
	- address-id-String-required: no-(filter condition) filter by the unique ID of the elastic public network IPv6.
	- public-ipv6-address-String-required: no-(filter condition) filter by the IP address of the public network IPv6.
	- charge-type-String-required: no-(filter condition) filter by billing type.
	- private-ipv6-address-String-required: no-(filter condition) filter by bound private network IPv6 address.
	- egress-String-required: no-(filter condition) filter by exit.
	- address-type-String-required: no-(filter condition) filter by IPv6 type.
	- address-isp-String-required: no-(filter condition) filter by operator type.
  The status includes: 'CREATING','BINDING','BIND','UNBINDING','UNBIND','OFFLINING','BIND_ENI','PRIVATE'.
	- address-name-String-required: no-(filter condition) filter by EIP name. Blur filtering is not supported.
	- tag-key-String-required: no-(filter condition) filter by label key.
	- tag-value-String-required: no-(filter condition) filter by tag value.
	- tag:tag-key-String-required: no-(filter condition) filter by label key value pair. Tag-key is replaced with a specific label key.
* `ipv6_address_ids` - (Optional, Set: [`String`]) Unique ID column that identifies IPv6.
	- Traditional Elastic IPv6 unique ID is like: `eip-11112222`
	- Elastic IPv6 unique ID is like: `eipv6 -11112222`
Note: Parameters do not support specifying both IPv6AddressIds and Filters.
* `result_output_file` - (Optional, String) Used to save results.
* `traditional` - (Optional, Bool) Whether to query traditional IPv6 address information.

The `filters` object supports the following:

* `name` - (Required, String) Property name. If there are multiple Filters, the relationship between Filters is a logical AND (AND) relationship.
* `values` - (Required, Set) Attribute value. If there are multiple Values in the same Filter, the relationship between Values under the same Filter is a logical OR relationship. When the value type is a Boolean type, the value can be directly taken to the string TRUE or FALSE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `address_set` - List of IPv6 details.


