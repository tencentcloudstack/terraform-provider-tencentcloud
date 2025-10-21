---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_network_interface_limit"
sidebar_current: "docs-tencentcloud-datasource-vpc_network_interface_limit"
description: |-
  Use this data source to query detailed information of vpc network_interface_limit
---

# tencentcloud_vpc_network_interface_limit

Use this data source to query detailed information of vpc network_interface_limit

## Example Usage

```hcl
data "tencentcloud_vpc_network_interface_limit" "network_interface_limit" {
  instance_id = "ins-cr2rfq78"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) ID of a CVM instance or ENI to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `eni_private_ip_address_quantity` - Quota of IP addresses that can be allocated to each standard-mounted ENI.
* `eni_quantity` - Quota of ENIs mounted to a CVM instance in a standard way.
* `extend_eni_private_ip_address_quantity` - Quota of IP addresses that can be allocated to each extension-mounted ENI.Note: this field may return `null`, indicating that no valid values can be obtained.
* `extend_eni_quantity` - Quota of ENIs mounted to a CVM instance as an extensionNote: this field may return `null`, indicating that no valid values can be obtained.
* `sub_eni_private_ip_address_quantity` - The quota of IPs that can be assigned to each relayed ENI.Note: This field may return `null`, indicating that no valid values can be obtained.
* `sub_eni_quantity` - The quota of relayed ENIsNote: This field may return `null`, indicating that no valid values can be obtained.


