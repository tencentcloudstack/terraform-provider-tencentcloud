---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_cluster_vpc_fw_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_cluster_vpc_fw_switch"
description: |-
  Provides a resource to create a CFW cluster vpc fw switch
---

# tencentcloud_cfw_cluster_vpc_fw_switch

Provides a resource to create a CFW cluster vpc fw switch

## Example Usage

### If switch_mode is 2

```hcl
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 2
  routing_mode = 0
  region_cidr_configs {
    region      = "ap-guangzhou"
    cidr_mode   = 1
    custom_cidr = ""
  }
}
```

### If switch_mode is 1

```hcl
resource "tencentcloud_cfw_cluster_vpc_fw_switch" "example" {
  ccn_id       = "ccn-8qv0ro89"
  switch_mode  = 1
  routing_mode = 1
  region_cidr_configs {
    region      = "ap-guangzhou"
    cidr_mode   = 0
    custom_cidr = ""
  }

  region_cidr_configs {
    region      = "ap-chongqing"
    cidr_mode   = 0
    custom_cidr = ""
  }

  region_cidr_configs {
    region      = "ap-shanghai"
    cidr_mode   = 1
    custom_cidr = ""
  }

  interconnect_pairs {
    interconnect_mode = "FullMesh"
    group_a {
      instance_id      = "vpc-264i7uzy"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.124.0.0/16"
      ]
    }

    group_a {
      instance_id      = "vpc-h2i9m8xh"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.25.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-264i7uzy"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.124.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-h2i9m8xh"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.25.0.0/16"
      ]
    }
  }

  interconnect_pairs {
    interconnect_mode = "CrossConnect"
    group_a {
      instance_id      = "vpc-5l5uqrgx"
      instance_type    = "VPC"
      instance_region  = "ap-chongqing"
      access_cidr_mode = 1
      access_cidr_list = [
        "192.168.0.0/16"
      ]
    }

    group_b {
      instance_id      = "vpc-1yoh1nhh"
      instance_type    = "VPC"
      instance_region  = "ap-guangzhou"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.208.0.0/24",
        "172.16.0.0/16"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN ID.
* `switch_mode` - (Required, Int, ForceNew) Switch access mode, 1: automatic access, 2: manual access.
* `interconnect_pairs` - (Optional, List) Interconnect pair list.
* `region_cidr_configs` - (Optional, List) Regional level CIDR configuration.
* `routing_mode` - (Optional, Int, ForceNew) Traffic steering routing method, 0: multi-route table, 1: policy routing.

The `group_a` object of `interconnect_pairs` supports the following:

* `access_cidr_list` - (Required, Set) List of network segments for accessing firewall.
* `access_cidr_mode` - (Required, Int) Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
* `instance_id` - (Required, String) Instance ID.
* `instance_region` - (Required, String) Region where the instance is located.
* `instance_type` - (Required, String) Instance type such as VPC or DIRECTCONNECT.

The `group_b` object of `interconnect_pairs` supports the following:

* `access_cidr_list` - (Required, Set) List of network segments for accessing firewall.
* `access_cidr_mode` - (Required, Int) Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
* `instance_id` - (Required, String) Instance ID.
* `instance_region` - (Required, String) Region where the instance is located.
* `instance_type` - (Required, String) Instance type such as VPC or DIRECTCONNECT.

The `interconnect_pairs` object supports the following:

* `group_a` - (Required, List) Group A.
* `group_b` - (Required, List) Group B.
* `interconnect_mode` - (Required, String) Interconnect mode: `CrossConnect`: cross interconnect (each instance in group A interconnects with each instance in group B), `FullMesh`: full mesh (group A content is identical to group B, equivalent to pairwise interconnection within the group).

The `region_cidr_configs` object supports the following:

* `cidr_mode` - (Required, Int) CIDR mode: 0-skip, 1-automatic, 2-custom.
* `custom_cidr` - (Required, String) Custom CIDR (required when CidrMode=2), empty string otherwise.
* `region` - (Required, String) Traffic steering region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeClusterVpcFwSwitchs](https://cloud.tencent.com/document/api/1132/127064) - 查询集群模式Vpc间防火墙开关
* Create: [DescribeSwitchError](https://cloud.tencent.com/document/api/1132/127053) - 互联网边界和VPC防火墙开关横幅错误信息
* Create: [ModifyClusterVpcFwSwitch](https://cloud.tencent.com/document/api/1132/127061) - 修改集群模式VPC防火墙开关
* Read: [DescribeCcnVpcFwSwitch](https://cloud.tencent.com/document/api/1132/127058) - 查询CCN VPC防火墙开关配置
* Read: [DescribeClusterVpcFwSwitchs](https://cloud.tencent.com/document/api/1132/127064) - 查询集群模式Vpc间防火墙开关
* Update: [DescribeClusterVpcFwSwitchs](https://cloud.tencent.com/document/api/1132/127064) - 查询集群模式Vpc间防火墙开关
* Update: [DescribeSwitchError](https://cloud.tencent.com/document/api/1132/127053) - 互联网边界和VPC防火墙开关横幅错误信息
* Update: [UpdateClusterVpcFw](https://cloud.tencent.com/document/api/1132/127055) - 修改更新CCN中VPC防火墙策略配置
* Delete: [DescribeClusterVpcFwSwitchs](https://cloud.tencent.com/document/api/1132/127064) - 查询集群模式Vpc间防火墙开关
* Delete: [DescribeSwitchError](https://cloud.tencent.com/document/api/1132/127053) - 互联网边界和VPC防火墙开关横幅错误信息
* Delete: [ModifyClusterVpcFwSwitch](https://cloud.tencent.com/document/api/1132/127061) - 修改集群模式VPC防火墙开关

