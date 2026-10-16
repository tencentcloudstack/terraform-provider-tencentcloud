---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_nat_instance"
sidebar_current: "docs-tencentcloud-resource-cfw_nat_instance"
description: |-
  Provides a resource to create a CFW nat instance
---

# tencentcloud_cfw_nat_instance

Provides a resource to create a CFW nat instance

## Example Usage

### If mode is 0

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name         = "tf_example"
  cross_a_zone = 1
  width        = 20
  mode         = 0
  new_mode_items {
    vpc_list = [
      "vpc-40hif9or"
    ]
    eips = [
      "119.29.107.37"
    ]
  }
  zone_set = [
    "ap-guangzhou-6",
    "ap-guangzhou-7"
  ]
}
```

### If mode is 1

```hcl
resource "tencentcloud_cfw_nat_instance" "example" {
  name         = "tf_example"
  cross_a_zone = 1
  width        = 20
  mode         = 1
  nat_gw_list = [
    "nat-9wwkz1kr"
  ]

  zone_set = [
    "ap-guangzhou-6",
    "ap-guangzhou-7"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `mode` - (Required, Int) Mode 1: access mode; 0: new mode.
* `name` - (Required, String) Firewall instance name.
* `width` - (Required, Int) Bandwidth.
* `zone_set` - (Required, Set: [`String`]) Zone list.
* `cross_a_zone` - (Optional, Int) Off-site disaster recovery 1: use off-site disaster recovery; 0: do not use off-site disaster recovery; if empty, the default is not to use off-site disaster recovery.
* `nat_gw_list` - (Optional, Set: [`String`]) A list of nat gateways connected to the access mode, at least one of NewModeItems and NatgwList is passed.
* `new_mode_items` - (Optional, List) New mode passing parameters are added, at least one of new_mode_items and nat_gw_list is passed.

The `new_mode_items` object supports the following:

* `eips` - (Required, Set) List of egress elastic public network IPs bound in the new mode.
* `vpc_list` - (Required, Set) List of vpcs connected in new mode.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `nat_instance_id` - Nat instance ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNatFwInstanceWithDomain](https://www.tencentcloud.com/search/?wd=CreateNatFwInstanceWithDomain&type=document&lang=en)
* Create: [DescribeNatFwInstancesInfo](https://www.tencentcloud.com/search/?wd=DescribeNatFwInstancesInfo&type=document&lang=en)
* Read: [DescribeCfwEips](https://www.tencentcloud.com/search/?wd=DescribeCfwEips&type=document&lang=en)
* Read: [DescribeNatFwInstancesInfo](https://www.tencentcloud.com/search/?wd=DescribeNatFwInstancesInfo&type=document&lang=en)
* Read: [DescribeNatFwVpcDnsLst](https://www.tencentcloud.com/search/?wd=DescribeNatFwVpcDnsLst&type=document&lang=en)
* Update: [ModifyNatInstance](https://www.tencentcloud.com/search/?wd=ModifyNatInstance&type=document&lang=en)
* Delete: [DeleteNatFwInstance](https://www.tencentcloud.com/search/?wd=DeleteNatFwInstance&type=document&lang=en)

## Import

CFW nat instance can be imported using the id, e.g.

```
terraform import tencentcloud_cfw_nat_instance.example cfwnat-54a21421
```

