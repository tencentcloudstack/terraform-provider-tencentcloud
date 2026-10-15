---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_hourdb_instance"
sidebar_current: "docs-tencentcloud-resource-dcdb_hourdb_instance"
description: |-
  Provides a resource to create a DCDB hourdb instance
---

# tencentcloud_dcdb_hourdb_instance

Provides a resource to create a DCDB hourdb instance

## Example Usage

```hcl
resource "tencentcloud_dcdb_hourdb_instance" "example" {
  instance_name     = "tf-example"
  zones             = ["ap-guangzhou-6", "ap-guangzhou-7"]
  shard_memory      = "4"
  shard_storage     = "50"
  shard_node_count  = "2"
  shard_count       = "2"
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  security_group_id = "sg-4z20n68d"
  db_version_id     = "8.0"
  resource_tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `shard_count` - (Required, Int) instance shard count.
* `shard_memory` - (Required, Int) memory(GB) for each shard. It can be obtained by querying api DescribeShardSpec.
* `shard_node_count` - (Required, Int) node count for each shard. It can be obtained by querying api DescribeShardSpec.
* `shard_storage` - (Required, Int) storage(GB) for each shard. It can be obtained by querying api DescribeShardSpec.
* `db_version_id` - (Optional, String) db engine version, default to Percona 5.7.17.
* `dcn_instance_id` - (Optional, String) DCN source instance ID.
* `dcn_region` - (Optional, String) DCN source region.
* `extranet_access` - (Optional, Bool) Whether to open the extranet access.
* `instance_name` - (Optional, String) name of this instance.
* `ipv6_flag` - (Optional, Int) Whether to support IPv6.
* `project_id` - (Optional, Int) project id.
* `resource_tags` - (Optional, List) resource tags.
* `security_group_id` - (Optional, String) security group id.
* `subnet_id` - (Optional, String) subnet id, its required when vpcId is set.
* `vip` - (Optional, String) The field is required to specify VIP.
* `vipv6` - (Optional, String) The field is required to specify VIPv6.
* `vpc_id` - (Optional, String) vpc id.
* `zones` - (Optional, Set: [`String`]) available zone.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) tag key.
* `tag_value` - (Required, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `vport` - Intranet port.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CloseDBExtranetAccess](https://cloud.tencent.com/document/api/557/19998) - 关闭外网访问
* Create: [CreateHourDCDBInstance](https://cloud.tencent.com/document/api/557/70249) - 创建TDSQL按量计费实例
* Create: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - 查询实例列表
* Create: [DescribeDcnDetail](https://cloud.tencent.com/document/api/557/56346) - 获取实例灾备详情
* Create: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - 查询流程状态
* Create: [InitDCDBInstances](https://cloud.tencent.com/document/api/557/19985) - 初始化实例
* Create: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/557/82021) - 修改实例所属网络
* Create: [OpenDBExtranetAccess](https://cloud.tencent.com/document/api/557/19980) - 开通外网访问
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/557/51781) - 查询实例安全组信息
* Read: [DescribeDCDBInstanceDetail](https://cloud.tencent.com/document/api/557/89385) - 获取实例详情
* Read: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - 查询实例列表
* Read: [DescribeDcnDetail](https://cloud.tencent.com/document/api/557/56346) - 获取实例灾备详情
* Update: [CloseDBExtranetAccess](https://cloud.tencent.com/document/api/557/19998) - 关闭外网访问
* Update: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - 查询实例列表
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - 查询流程状态
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/557/73095) - 修改实例名字
* Update: [ModifyDBInstancesProject](https://cloud.tencent.com/document/api/557/19983) - 修改实例所属项目
* Update: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/557/82021) - 修改实例所属网络
* Update: [OpenDBExtranetAccess](https://cloud.tencent.com/document/api/557/19980) - 开通外网访问
* Update: [UpgradeHourDCDBInstance](https://cloud.tencent.com/document/api/557/85926) - 升级TDSQL按量计费实例
* Delete: [DestroyHourDCDBInstance](https://cloud.tencent.com/document/api/557/56483) - 销毁TDSQL按量计费实例

## Import

DCDB hourdb instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.example tdsqlshard-nr6j5sed
```

