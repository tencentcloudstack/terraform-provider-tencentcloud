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

* Create: [CloseDBExtranetAccess](https://www.tencentcloud.com/search/?wd=CloseDBExtranetAccess&type=document&lang=en)
* Create: [CreateHourDCDBInstance](https://www.tencentcloud.com/search/?wd=CreateHourDCDBInstance&type=document&lang=en)
* Create: [DescribeDCDBInstances](https://www.tencentcloud.com/search/?wd=DescribeDCDBInstances&type=document&lang=en)
* Create: [DescribeDcnDetail](https://www.tencentcloud.com/search/?wd=DescribeDcnDetail&type=document&lang=en)
* Create: [DescribeFlow](https://www.tencentcloud.com/search/?wd=DescribeFlow&type=document&lang=en)
* Create: [InitDCDBInstances](https://www.tencentcloud.com/search/?wd=InitDCDBInstances&type=document&lang=en)
* Create: [ModifyInstanceNetwork](https://www.tencentcloud.com/search/?wd=ModifyInstanceNetwork&type=document&lang=en)
* Create: [OpenDBExtranetAccess](https://www.tencentcloud.com/search/?wd=OpenDBExtranetAccess&type=document&lang=en)
* Read: [DescribeDBSecurityGroups](https://www.tencentcloud.com/search/?wd=DescribeDBSecurityGroups&type=document&lang=en)
* Read: [DescribeDCDBInstanceDetail](https://www.tencentcloud.com/search/?wd=DescribeDCDBInstanceDetail&type=document&lang=en)
* Read: [DescribeDCDBInstances](https://www.tencentcloud.com/search/?wd=DescribeDCDBInstances&type=document&lang=en)
* Read: [DescribeDcnDetail](https://www.tencentcloud.com/search/?wd=DescribeDcnDetail&type=document&lang=en)
* Update: [CloseDBExtranetAccess](https://www.tencentcloud.com/search/?wd=CloseDBExtranetAccess&type=document&lang=en)
* Update: [DescribeDCDBInstances](https://www.tencentcloud.com/search/?wd=DescribeDCDBInstances&type=document&lang=en)
* Update: [DescribeFlow](https://www.tencentcloud.com/search/?wd=DescribeFlow&type=document&lang=en)
* Update: [ModifyDBInstanceName](https://www.tencentcloud.com/search/?wd=ModifyDBInstanceName&type=document&lang=en)
* Update: [ModifyDBInstancesProject](https://www.tencentcloud.com/search/?wd=ModifyDBInstancesProject&type=document&lang=en)
* Update: [ModifyInstanceNetwork](https://www.tencentcloud.com/search/?wd=ModifyInstanceNetwork&type=document&lang=en)
* Update: [OpenDBExtranetAccess](https://www.tencentcloud.com/search/?wd=OpenDBExtranetAccess&type=document&lang=en)
* Update: [UpgradeHourDCDBInstance](https://www.tencentcloud.com/search/?wd=UpgradeHourDCDBInstance&type=document&lang=en)
* Delete: [DestroyHourDCDBInstance](https://www.tencentcloud.com/search/?wd=DestroyHourDCDBInstance&type=document&lang=en)

## Import

DCDB hourdb instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.example tdsqlshard-nr6j5sed
```

