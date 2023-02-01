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
* `shard_memory` - (Required, Int) &amp;quot;Shard memory size, unit: GB, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.
* `shard_node_count` - (Required, Int) &amp;quot;Number of single shard nodes, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.
* `shard_storage` - (Required, Int) &amp;quot;Shard storage size, unit: GB, can pass DescribeShardSpec&amp;quot;&amp;quot;Query the instance specification to obtain.&amp;quot;.
* `zones` - (Required, Set: [`String`]) &amp;quot;The availability zone distribution of shard nodes can be filled with up to two availability zones. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.&amp;quot;&amp;quot;Note that the current availability zone that can be sold needs to be pulled through the DescribeDCDBSaleInfo interface.&amp;quot;.
* `auto_renew_flag` - (Optional, Int) &amp;quot;Automatic renewal flag, 0 means the default state (the user has not set it, that is, the initial state is manual renewal, and the user has activated the prepaid non-stop privilege and will also perform automatic renewal).&amp;quot;&amp;quot;1 means automatic renewal, 2 means no automatic renewal (user setting).&amp;quot;&amp;quot;if the business has no concept of renewal or automatic renewal is not required, it needs to be set to 0.&amp;quot;.
* `auto_voucher` - (Optional, Bool) Whether to automatically use vouchers for payment, not used by default.
* `db_version_id` - (Optional, String) &amp;quot;Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17.&amp;quot;&amp;quot;8.0.18 - MySQL 8.0.18;&amp;quot;&amp;quot;10.1.9 - Mariadb 10.1.9;&amp;quot;&amp;quot;5.7.17 - Percona 5.7.17&amp;quot;&amp;quot;If not filled, the default is 5.7.17, which means Percona 5.7.17.&amp;quot;.
* `dcn_instance_id` - (Optional, String) DCN source instance ID.
* `dcn_region` - (Optional, String) DCN source region.
* `init_params` - (Optional, List) &amp;quot;parameter list. The optional values of this interface are:&amp;quot;&amp;quot;character_set_server (character set, must be passed),&amp;quot;&amp;quot;lower_case_table_names (table name is case sensitive, must be passed, 0 - sensitive; 1 - insensitive),&amp;quot;&amp;quot;innodb_page_size (innodb data page, default 16K),&amp;quot;&amp;quot;sync_mode ( Synchronous mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous degenerate. The default is strong synchronous degenerate)&amp;quot;.
* `instance_name` - (Optional, String) Instance name, you can set the name of the instance independently through this field.
* `ipv6_flag` - (Optional, Int) Whether to support IPv6.
* `project_id` - (Optional, Int) Project ID, which can be obtained by viewing the project list, if not passed, it will be associated with the default project.
* `resource_tags` - (Optional, List) Array of tag key-value pairs.
* `security_group_ids` - (Optional, Set: [`String`]) Security group ids, the security group can be passed in the form of an array, compatible with the previous SecurityGroupId parameter.
* `subnet_id` - (Optional, String) Virtual private network subnet ID, required when VpcId is not empty.
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



## Import

dcdb db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_db_instance.db_instance db_instance_id
```

