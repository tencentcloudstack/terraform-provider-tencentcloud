---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnats"
sidebar_current: "docs-tencentcloud-datasource-dnats"
description: |-
  Use this data source to query detailed information of DNATs.
---

# tencentcloud_dnats

Use this data source to query detailed information of DNATs.

## Example Usage

```hcl
# query by nat gateway id
data "tencentcloud_dnats" "foo" {
  nat_id = "nat-xfaq1"
}

# query by vpc id
data "tencentcloud_dnats" "foo" {
  vpc_id = "vpc-xfqag"
}

# query by elastic ip
data "tencentcloud_dnats" "foo" {
  elastic_ip = "123.207.115.136"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the NAT forward.
* `elastic_ip` - (Optional) Network address of the EIP.
* `elastic_port` - (Optional) Port of the EIP.
* `nat_id` - (Optional) Id of the NAT gateway.
* `private_ip` - (Optional) Network address of the backend service.
* `private_port` - (Optional) Port of intranet.
* `result_output_file` - (Optional) Used to save results.
* `vpc_id` - (Optional) Id of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dnat_list` - Information list of the DNATs.
  * `elastic_ip` - Network address of the EIP.
  * `elastic_port` - Port of the EIP.
  * `nat_id` - Id of the NAT.
  * `private_ip` - Network address of the backend service.
  * `private_port` - Port of intranet.
  * `protocol` - Type of the network protocol. Valid values: `TCP` and `UDP`.
  * `vpc_id` - Id of the VPC.


