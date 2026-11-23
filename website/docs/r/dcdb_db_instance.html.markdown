---
subcategory: "TDSQL for MySQL(DCDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dcdb_db_instance"
sidebar_current: "docs-tencentcloud-resource-dcdb_db_instance"
description: |-
  Provides a resource to create a dcdb db_instance
---

# tencentcloud_dcdb_db_instance

Provides a resource to create a dcdb db_instance

## Example Usage

```hcl
resource "tencentcloud_dcdb_db_instance" "db_instance" {
  instance_name    = "test_dcdb_db_instance"
  zones            = ["ap-guangzhou-5"]
  period           = 1
  shard_memory     = "2"
  shard_storage    = "10"
  shard_node_count = "2"
  shard_count      = "2"
  vpc_id           = local.vpc_id
  subnet_id        = local.subnet_id
  db_version_id    = "8.0"
  resource_tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "1"
  }
  init_params {
    param = "sync_mode"
    value = "2"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  security_group_ids = [local.sg_id]
}
```

## Argument Reference

The following arguments are supported:

* `period` - (Required, Int) The length of time you want to buy, unit: month.
* `shard_count` - (Required, Int) The number of instance fragments, the optional range is 2-8, and new fragments can be added to a maximum of 64 fragments by upgrading the instance.
* `shard_memory` - (Required, Int) Shard memory size, unit: GB, can pass DescribeShardSpec Query the instance specification to obtain.
* `shard_node_count` - (Required, Int) Number of single shard nodes, can pass DescribeShardSpec  Query the instance specification to obtain.
* `shard_storage` - (Required, Int) Shard storage size, unit: GB, can pass DescribeShardSpec  Query the instance specification to obtain.
* `zones` - (Required, Set: [`String`]) The availability zone distribution of shard nodes can be filled with up to two availability zones. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.Note that the current availability zone that can be sold needs to be pulled through the DescribeDCDBSaleInfo interface.
* `auto_renew_flag` - (Optional, Int) Automatic renewal flag, 0 means the default state (the user has not set it, that is, the initial state is manual renewal, and the user has activated the prepaid non-stop privilege and will also perform automatic renewal).  1 means automatic renewal, 2 means no automatic renewal (user setting).  if the business has no concept of renewal or automatic renewal is not required, it needs to be set to 0.
* `auto_voucher` - (Optional, Bool) Whether to automatically use vouchers for payment, not used by default.
* `db_version_id` - (Optional, String) Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17.  8.0.18 - MySQL 8.0.18;  10.1.9 - Mariadb 10.1.9;  5.7.17 - Percona 5.7.17  If not filled, the default is 5.7.17, which means Percona 5.7.17.
* `dcn_instance_id` - (Optional, String) DCN source instance ID.
* `dcn_region` - (Optional, String) DCN source region.
* `extranet_access` - (Optional, Bool) Whether to open the extranet access.
* `init_params` - (Optional, List) Parameter list. The optional values of this interface are:  character_set_server (character set, must be passed),  lower_case_table_names (table name is case sensitive, must be passed, 0 - sensitive; 1 - insensitive),  innodb_page_size (innodb data page, default 16K),  sync_mode (Synchronous mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous degenerate. The default is strong synchronous degenerate).
* `instance_name` - (Optional, String) Instance name, you can set the name of the instance independently through this field.
* `ipv6_flag` - (Optional, Int) Whether to support IPv6.
* `project_id` - (Optional, Int) Project ID, which can be obtained by viewing the project list, if not passed, it will be associated with the default project.
* `resource_tags` - (Optional, List) Array of tag key-value pairs.
* `security_group_ids` - (Optional, Set: [`String`]) Security group ids, the security group can be passed in the form of an array, compatible with the previous SecurityGroupId parameter.
* `subnet_id` - (Optional, String) Virtual private network subnet ID, required when VpcId is not empty.
* `vip` - (Optional, String) The field is required to specify VIP.
* `vipv6` - (Optional, String) The field is required to specify VIPv6.
* `voucher_ids` - (Optional, Set: [`String`]) Voucher ID list, currently only supports specifying one voucher.
* `vpc_id` - (Optional, String) Virtual private network ID, if not passed or passed empty, it means that it is created as a basic network.

The `init_params` object supports the following:

* `param` - (Required, String) The name of parameter.
* `value` - (Required, String) The value of parameter.

The `resource_tags` object supports the following:

* `tag_key` - (Required, String) The key of tag.
* `tag_value` - (Required, String) The value of tag.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `vport` - Intranet port.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [CloseDBExtranetAccess](https://www.tencentcloud.com/document/api/1042/34451)
  - [CreateDCDBInstance](https://www.tencentcloud.com/document/api/1042/52231)
  - [DescribeDCDBInstances](https://www.tencentcloud.com/document/api/1042/34442)
  - [DescribeDcnDetail](https://www.tencentcloud.com/document/api/1042/40671)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1042/40887)
  - [InitDCDBInstances](https://www.tencentcloud.com/document/api/1042/34436)
  - [ModifyInstanceNetwork](https://www.tencentcloud.com/document/api/1042/51398)
  - [OpenDBExtranetAccess](https://www.tencentcloud.com/document/api/1042/34431)
* Read:
  - [DescribeDBSecurityGroups](https://www.tencentcloud.com/document/api/1042/39459)
  - [DescribeDCDBInstanceDetail](https://www.tencentcloud.com/document/api/1042/54586)
  - [DescribeDCDBInstances](https://www.tencentcloud.com/document/api/1042/34442)
  - [DescribeDcnDetail](https://www.tencentcloud.com/document/api/1042/40671)
* Update:
  - [CloseDBExtranetAccess](https://www.tencentcloud.com/document/api/1042/34451)
  - [DescribeFlow](https://www.tencentcloud.com/document/api/1042/40887)
  - [ModifyDBInstanceName](https://www.tencentcloud.com/document/api/1042/46993)
  - [ModifyDBInstancesProject](https://www.tencentcloud.com/document/api/1042/34434)
  - [ModifyInstanceNetwork](https://www.tencentcloud.com/document/api/1042/51398)
  - [OpenDBExtranetAccess](https://www.tencentcloud.com/document/api/1042/34431)
  - [RenewDCDBInstance](https://www.tencentcloud.com/search/?wd=RenewDCDBInstance&type=document&lang=en)
* Delete: [DestroyDCDBInstance](https://www.tencentcloud.com/document/api/1042/40886)

## Import

dcdb db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_instance.db_instance db_instance_id
```

