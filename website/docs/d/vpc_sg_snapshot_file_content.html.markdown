---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_sg_snapshot_file_content"
sidebar_current: "docs-tencentcloud-datasource-vpc_sg_snapshot_file_content"
description: |-
  Use this data source to query detailed information of vpc sg_snapshot_file_content
---

# tencentcloud_vpc_sg_snapshot_file_content

Use this data source to query detailed information of vpc sg_snapshot_file_content

## Example Usage

```hcl
data "tencentcloud_vpc_sg_snapshot_file_content" "sg_snapshot_file_content" {
  snapshot_policy_id = "sspolicy-ebjofe71"
  snapshot_file_id   = "ssfile-017gepjxpr"
  security_group_id  = "sg-ntrgm89v"
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required, String) Security group ID.
* `snapshot_file_id` - (Required, String) Snapshot file ID.
* `snapshot_policy_id` - (Required, String) Snapshot policy IDs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backup_data` - Backup data.
  * `action` - ACCEPT or DROP.
  * `address_template` - IP address ID or IP address group ID.
    * `address_group_id` - The ID of the IP address group, such as `ipmg-2uw6ujo6`.
    * `address_id` - The ID of the IP address, such as `ipm-2uw6ujo6`.
  * `cidr_block` - Either `CidrBlock` or `Ipv6CidrBlock can be specified. Note that if `0.0.0.0/n` is entered, it is mapped to 0.0.0.0/0.
  * `ipv6_cidr_block` - The CIDR block or IPv6 (mutually exclusive).
  * `modify_time` - The last modification time of the security group.
  * `policy_description` - Security group policy description.
  * `policy_index` - The index number of security group rules, which dynamically changes with the rules. This parameter can be obtained via the `DescribeSecurityGroupPolicies` API and used with the `Version` field in the returned value of the API.
  * `port` - Port (`all`, a single port, or a port range).Note: If the `Protocol` value is set to `ALL`, the `Port` value also needs to be set to `all`.
  * `protocol` - Protocol. Valid values: TCP, UDP, ICMP, ICMPv6, ALL.
  * `security_group_id` - The security group instance ID, such as `sg-ohuuioma`.
  * `service_template` - Protocol port ID or protocol port group ID. ServiceTemplate and Protocol+Port are mutually exclusive.
    * `service_group_id` - Protocol port group ID, such as `ppmg-f5n1f8da`.
    * `service_id` - Protocol port ID, such as `ppm-f5n1f8da`.
* `backup_time` - Backup time.
* `instance_id` - Security group ID.
* `operator` - Operator.
* `original_data` - Original data.
  * `action` - ACCEPT or DROP.
  * `address_template` - IP address ID or IP address group ID.
    * `address_group_id` - The ID of the IP address group, such as `ipmg-2uw6ujo6`.
    * `address_id` - The ID of the IP address, such as `ipm-2uw6ujo6`.
  * `cidr_block` - Either `CidrBlock` or `Ipv6CidrBlock can be specified. Note that if `0.0.0.0/n` is entered, it is mapped to 0.0.0.0/0.
  * `ipv6_cidr_block` - The CIDR block or IPv6 (mutually exclusive).
  * `modify_time` - The last modification time of the security group.
  * `policy_description` - Security group policy description.
  * `policy_index` - The index number of security group rules, which dynamically changes with the rules. This parameter can be obtained via the `DescribeSecurityGroupPolicies` API and used with the `Version` field in the returned value of the API.
  * `port` - Port (`all`, a single port, or a port range).Note: If the `Protocol` value is set to `ALL`, the `Port` value also needs to be set to `all`.
  * `protocol` - Protocol. Valid values: TCP, UDP, ICMP, ICMPv6, ALL.
  * `security_group_id` - The security group instance ID, such as `sg-ohuuioma`.
  * `service_template` - Protocol port ID or protocol port group ID. ServiceTemplate and Protocol+Port are mutually exclusive.
    * `service_group_id` - Protocol port group ID, such as `ppmg-f5n1f8da`.
    * `service_id` - Protocol port ID, such as `ppm-f5n1f8da`.


