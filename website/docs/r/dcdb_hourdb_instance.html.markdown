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

* Create: [CloseDBExtranetAccess](https://cloud.tencent.com/document/api/557/19998) - This API is used to disable public network access for a TencentDB instance, which will make the public IP address inaccessible. The `DescribeDCDBInstances` API will not return the public domain name and port information of the corresponding instance.
* Create: [CreateHourDCDBInstance](https://cloud.tencent.com/document/api/557/70249) - This API is used to create a pay-as-you-go TDSQL instance.
* Create: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - This API is used to query the list of TencentDB instances. It supports filtering instances by project ID, instance ID, private network address, and instance name.
* Create: [DescribeDcnDetail](https://cloud.tencent.com/document/api/557/56346) - This API is used to query the disaster recovery details of an instance.
* Create: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - This API is used to query task status.
* Create: [InitDCDBInstances](https://cloud.tencent.com/document/api/557/19985) - This API is used to initialize instances, including setting the default character set and table name case sensitivity.
* Create: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/557/82021) - This API is used to modify instance network.
* Create: [OpenDBExtranetAccess](https://cloud.tencent.com/document/api/557/19980)
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/557/51781) - This API is used to query the security group information of an instance.
* Read: [DescribeDCDBInstanceDetail](https://cloud.tencent.com/document/api/557/89385) - This API is used to get the details of a TDSQL instance.
* Read: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - This API is used to query the list of TencentDB instances. It supports filtering instances by project ID, instance ID, private network address, and instance name.
* Read: [DescribeDcnDetail](https://cloud.tencent.com/document/api/557/56346) - This API is used to query the disaster recovery details of an instance.
* Update: [CloseDBExtranetAccess](https://cloud.tencent.com/document/api/557/19998) - This API is used to disable public network access for a TencentDB instance, which will make the public IP address inaccessible. The `DescribeDCDBInstances` API will not return the public domain name and port information of the corresponding instance.
* Update: [DescribeDCDBInstances](https://cloud.tencent.com/document/api/557/16140) - This API is used to query the list of TencentDB instances. It supports filtering instances by project ID, instance ID, private network address, and instance name.
* Update: [DescribeFlow](https://cloud.tencent.com/document/api/557/56485) - This API is used to query task status.
* Update: [ModifyDBInstanceName](https://cloud.tencent.com/document/api/557/73095) - This API is used to modify instance name.
* Update: [ModifyDBInstancesProject](https://cloud.tencent.com/document/api/557/19983) - This API is used to modify the project to which TencentDB instances belong.
* Update: [ModifyInstanceNetwork](https://cloud.tencent.com/document/api/557/82021) - This API is used to modify instance network.
* Update: [OpenDBExtranetAccess](https://cloud.tencent.com/document/api/557/19980)
* Update: [UpgradeHourDCDBInstance](https://cloud.tencent.com/document/api/557/85926) - This API is used to upgrade a pay-as-you-go TDSQL instance.
* Delete: [DestroyHourDCDBInstance](https://cloud.tencent.com/document/api/557/56483) - This API is used to terminate a pay-as-you-go TDSQL instance.

## Import

DCDB hourdb instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.example tdsqlshard-nr6j5sed
```

