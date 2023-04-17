---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_assist_vpc"
sidebar_current: "docs-tencentcloud-resource-cvm_chc_assist_vpc"
description: |-
  Provides a resource to create a cvm chc_assist_vpc
---

# tencentcloud_cvm_chc_assist_vpc

Provides a resource to create a cvm chc_assist_vpc

## Example Usage

```hcl
resource "tencentcloud_cvm_chc_assist_vpc" "chc_assist_vpc" {
  chc_id = "chc-xxxxx"
  bmc_virtual_private_cloud {
    vpc_id    = "vpc-xxxxx"
    subnet_id = "subnet-xxxxx"

  }
  bmc_security_group_ids = ["sg-xxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `bmc_virtual_private_cloud` - (Required, List, ForceNew) Out-of-band network information.
* `chc_id` - (Required, String, ForceNew) CHC host ID.
* `bmc_security_group_ids` - (Optional, Set: [`String`], ForceNew) Out-of-band network security group list.
* `deploy_security_group_ids` - (Optional, Set: [`String`], ForceNew) Deployment network security group list.
* `deploy_virtual_private_cloud` - (Optional, List, ForceNew) Deployment network information.

The `bmc_virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String) VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
* `vpc_id` - (Required, String) VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `as_vpc_gateway` - (Optional, Bool) Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
* `ipv6_address_count` - (Optional, Int) Number of IPv6 addresses randomly generated for the ENI.
* `private_ip_addresses` - (Optional, Set) Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.

The `deploy_virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String) VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
* `vpc_id` - (Required, String) VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `as_vpc_gateway` - (Optional, Bool) Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
* `ipv6_address_count` - (Optional, Int) Number of IPv6 addresses randomly generated for the ENI.
* `private_ip_addresses` - (Optional, Set) Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cvm chc_assist_vpc can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_assist_vpc.chc_assist_vpc chc_assist_vpc_id
```

