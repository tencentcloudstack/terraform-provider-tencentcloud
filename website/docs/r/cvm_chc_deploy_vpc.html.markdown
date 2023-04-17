---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_deploy_vpc"
sidebar_current: "docs-tencentcloud-resource-cvm_chc_deploy_vpc"
description: |-
  Provides a resource to create a cvm chc_deploy_vpc
---

# tencentcloud_cvm_chc_deploy_vpc

Provides a resource to create a cvm chc_deploy_vpc

## Example Usage

```hcl
resource "tencentcloud_cvm_chc_deploy_vpc" "chc_deploy_vpc" {
  chc_id = "chc-xxxxx"
  deploy_virtual_private_cloud {
    vpc_id    = "vpc-xxxxx"
    subnet_id = "subnet-xxxxx"
  }
  deploy_security_group_ids = ["sg-xxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `chc_id` - (Required, String, ForceNew) CHC host ID.
* `deploy_virtual_private_cloud` - (Required, List, ForceNew) Deployment network information.
* `deploy_security_group_ids` - (Optional, Set: [`String`], ForceNew) Deployment network security group list.

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

cvm chc_deploy_vpc can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_deploy_vpc.chc_deploy_vpc chc_deploy_vpc_id
```

