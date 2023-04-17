---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_config"
sidebar_current: "docs-tencentcloud-resource-cvm_chc_config"
description: |-
  Provides a resource to create a cvm chc_config
---

# tencentcloud_cvm_chc_config

Provides a resource to create a cvm chc_config

## Example Usage

```hcl
resource "tencentcloud_cvm_chc_config" "chc_config" {
  chc_id        = "chc-xxxxxx"
  instance_name = "xxxxxx"
  bmc_user      = "admin"
  password      = "xxxxxx"
  bmc_virtual_private_cloud {
    vpc_id    = "vpc-xxxxxx"
    subnet_id = "subnet-xxxxxx"

  }
  bmc_security_group_ids = ["sg-xxxxxx"]

  deploy_virtual_private_cloud {
    vpc_id    = "vpc-xxxxxx"
    subnet_id = "subnet-xxxxxx"
  }
  deploy_security_group_ids = ["sg-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `chc_id` - (Required, String, ForceNew) CHC host ID.
* `bmc_security_group_ids` - (Optional, List: [`String`], ForceNew) Out-of-band network security group list.
* `bmc_user` - (Optional, String) Valid characters: Letters, numbers, hyphens and underscores. Only set when update password.
* `bmc_virtual_private_cloud` - (Optional, List, ForceNew) Out-of-band network information.
* `deploy_security_group_ids` - (Optional, List: [`String`], ForceNew) Deployment network security group list.
* `deploy_virtual_private_cloud` - (Optional, List, ForceNew) Deployment network information.
* `device_type` - (Optional, String) Server type.
* `instance_name` - (Optional, String) CHC host name.
* `password` - (Optional, String) The password can contain 8 to 16 characters, including letters, numbers and special symbols (()`~!@#$%^&amp;amp;*-+=_|{}).

The `bmc_virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String, ForceNew) VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
* `vpc_id` - (Required, String, ForceNew) VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `as_vpc_gateway` - (Optional, Bool, ForceNew) Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
* `ipv6_address_count` - (Optional, Int, ForceNew) Number of IPv6 addresses randomly generated for the ENI.
* `private_ip_addresses` - (Optional, List, ForceNew) Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.

The `deploy_virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String, ForceNew) VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
* `vpc_id` - (Required, String, ForceNew) VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `as_vpc_gateway` - (Optional, Bool, ForceNew) Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
* `ipv6_address_count` - (Optional, Int, ForceNew) Number of IPv6 addresses randomly generated for the ENI.
* `private_ip_addresses` - (Optional, List, ForceNew) Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cvm chc_config can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_chc_config.chc_config chc_config_id
```

