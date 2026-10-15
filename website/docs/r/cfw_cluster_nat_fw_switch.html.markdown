---
subcategory: "Cloud Firewall(CFW)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfw_cluster_nat_fw_switch"
sidebar_current: "docs-tencentcloud-resource-cfw_cluster_nat_fw_switch"
description: |-
  Provides a resource to create a CFW cluster NAT firewall switch
---

# tencentcloud_cfw_cluster_nat_fw_switch

Provides a resource to create a CFW cluster NAT firewall switch

## Example Usage

```hcl
resource "tencentcloud_cfw_cluster_nat_fw_switch" "example" {
  nat_ccn_switch {
    nat_ins_id   = "nat-h1i1mf4n"
    ccn_id       = "ccn-p3mlp0tj"
    switch_mode  = 1
    routing_mode = 1
    access_instance_list {
      instance_id      = "vpc-8za5vlpc"
      instance_type    = "VPC"
      instance_region  = "ap-shanghai"
      access_cidr_mode = 1
      access_cidr_list = [
        "10.51.0.0/16"
      ]
    }

    access_instance_list {
      instance_id      = "vpc-ec6krnpp"
      instance_type    = "VPC"
      instance_region  = "ap-guangzhou"
      access_cidr_mode = 1
      access_cidr_list = [
        "172.17.0.0/16"
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `nat_ccn_switch` - (Required, List) NAT CCN switch configuration.

The `access_instance_list` object of `nat_ccn_switch` supports the following:

* `access_cidr_mode` - (Required, Int) Network segment mode for accessing firewall: 0-no access, 1-access all network segments associated with the instance, 2-access user-defined network segments.
* `instance_id` - (Required, String) Instance ID.
* `instance_region` - (Required, String) Region where the instance is located.
* `instance_type` - (Required, String) Instance type such as VPC or DIRECTCONNECT.
* `access_cidr_list` - (Optional, List) List of network segments for accessing firewall.

The `nat_ccn_switch` object supports the following:

* `ccn_id` - (Required, String, ForceNew) CCN instance ID.
* `nat_ins_id` - (Required, String, ForceNew) NAT firewall instance ID.
* `switch_mode` - (Required, Int) Switch access mode, 1: automatic access, 2: manual access.
* `access_instance_list` - (Optional, List) List of access instances.
* `lead_vpc_cidr` - (Optional, String) CIDR of the lead VPC.
* `routing_mode` - (Optional, Int) Traffic steering routing method, 0: multi-route table, 1: policy routing. Automatic access mode only supports policy routing (1); manual access mode supports both multi-route table (0) and policy routing (1).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeClusterNatCcnFwSwitchList](https://cloud.tencent.com/document/api/1132/132581) - 查询NAT CCN集群模式防火墙开关列表
* Create: [OpenClusterNatFwSwitch](https://cloud.tencent.com/document/api/1132/132577) - 开启NAT CCN集群模式防火墙开关
* Read: [DescribeNatCcnFwSwitch](https://cloud.tencent.com/document/api/1132/132580) - 查询NAT CCN防火墙开关配置
* Update: [DescribeClusterNatCcnFwSwitchList](https://cloud.tencent.com/document/api/1132/132581) - 查询NAT CCN集群模式防火墙开关列表
* Update: [ModifyClusterNatFwSwitch](https://cloud.tencent.com/document/api/1132/132578) - 修改NAT CCN集群模式防火墙开关配置
* Delete: [CloseClusterNatFwSwitch](https://cloud.tencent.com/document/api/1132/132582) - 关闭NAT CCN集群模式防火墙开关
* Delete: [DescribeClusterNatCcnFwSwitchList](https://cloud.tencent.com/document/api/1132/132581) - 查询NAT CCN集群模式防火墙开关列表

## Import

CFW cluster NAT firewall switch can be imported using the composite id `nat_ins_id#ccn_id`, e.g.

```
terraform import tencentcloud_cfw_cluster_nat_fw_switch.example nat-h1i1mf4n#ccn-p3mlp0tj
```

