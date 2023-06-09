---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_instance"
sidebar_current: "docs-tencentcloud-resource-mariadb_instance"
description: |-
  Provides a resource to create a mariadb instance
---

# tencentcloud_mariadb_instance

Provides a resource to create a mariadb instance

## Example Usage

```hcl
resource "tencentcloud_mariadb_instance" "instance" {
  zones      = ["ap-guangzhou-3", ]
  node_count = 2
  memory     = 8
  storage    = 10
  period     = 1
  # auto_voucher =
  # voucher_ids =
  vpc_id    = "vpc-ii1jfbhl"
  subnet_id = "subnet-3ku415by"
  # project_id = ""
  db_version_id = "8.0"
  instance_name = "terraform-test"
  # security_group_ids = ""
  auto_renew_flag = 1
  ipv6_flag       = 0
  tags = {
    "createby" = "terrafrom-2"
  }
  init_params {
    param = "character_set_server"
    value = "utf8mb4"
  }
  init_params {
    param = "lower_case_table_names"
    value = "0"
  }
  init_params {
    param = "innodb_page_size"
    value = "16384"
  }
  init_params {
    param = "sync_mode"
    value = "1"
  }
  dcn_region      = ""
  dcn_instance_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required, Int) Memory size, unit: GB, can be obtained by querying instance specifications through DescribeDBInstanceSpecs.
* `node_count` - (Required, Int) Number of nodes, 2 is one master and one slave, 3 is one master and two slaves.
* `storage` - (Required, Int) Storage size, unit: GB. You can query instance specifications through DescribeDBInstanceSpecs to obtain the lower and upper limits of disk specifications corresponding to different memory sizes.
* `zones` - (Required, Set: [`String`]) Instance node availability zone distribution, up to two availability zones can be filled. When the shard specification is one master and two slaves, two of the nodes are in the first availability zone.
* `auto_renew_flag` - (Optional, Int) Automatic renewal flag, 1: automatic renewal, 2: no automatic renewal.
* `auto_voucher` - (Optional, Bool, ForceNew) Whether to automatically use the voucher for payment, the default is not used.
* `db_version_id` - (Optional, String) Database engine version, currently available: 8.0.18, 10.1.9, 5.7.17. If not passed, the default is Percona 5.7.17.
* `dcn_instance_id` - (Optional, String, ForceNew) DCN source instance ID.
* `dcn_region` - (Optional, String, ForceNew) DCN source region.
* `init_params` - (Optional, List, ForceNew) Parameter list. The optional values of this interface are: character_set_server (character set, required) enum: utf8,latin1,gbk,utf8mb4,gb18030, lower_case_table_names (table name case sensitive, required, 0 - sensitive; 1 - insensitive), innodb_page_size (innodb data page, Default 16K), sync_mode (sync mode: 0 - asynchronous; 1 - strong synchronous; 2 - strong synchronous can degenerate. The default is strong synchronous can degenerate).
* `instance_name` - (Optional, String) Instance name, you can set the name of the instance independently through this field.
* `ipv6_flag` - (Optional, Int) Whether IPv6 is supported.
* `period` - (Optional, Int, ForceNew) The duration of the purchase, unit: month.
* `project_id` - (Optional, Int) Project ID, which can be obtained by viewing the project list, if not passed, it will be associated with the default project.
* `security_group_ids` - (Optional, Set: [`String`]) Security group ID list.
* `subnet_id` - (Optional, String) Virtual private network subnet ID, required when VpcId is not empty.
* `tags` - (Optional, Map) tag list.
* `vip` - (Optional, String) Intranet IP address.
* `voucher_ids` - (Optional, Set: [`String`], ForceNew) A list of voucher IDs. Currently, only one voucher can be specified.
* `vpc_id` - (Optional, String) Virtual private network ID, if not passed, it means that it is created as a basic network.

The `init_params` object supports the following:

* `param` - (Required, String) parameter name.
* `value` - (Required, String) parameter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - ID of the application to which the instance belongs.
* `cpu` - Number of CPU cores of the instance.
* `create_time` - Instance creation time, the format is 2006-01-02 15:04:05.
* `db_engine` - Database Engine.
* `dcn_dst_num` - Number of DCN disaster recovery instances.
* `dcn_flag` - DCN flag, 0-none, 1-primary instance, 2-disaster backup instance.
* `dcn_status` - DCN status, 0-none, 1-creating, 2-synchronizing, 3-disconnected.
* `excluster_id` - Exclusive cluster ID, if it is empty, it means a normal instance.
* `instance_id` - Instance ID, uniquely identifies a TDSQL instance.
* `instance_type` - 1: primary instance (exclusive), 2: primary instance, 3: disaster recovery instance, 4: disaster recovery instance (exclusive type).
* `is_audit_supported` - Whether the instance supports auditing. 1-supported; 0-not supported.
* `is_encrypt_supported` - Whether data encryption is supported. 1-supported; 0-not supported.
* `is_tmp` - Whether it is a temporary instance, 0 means no, non-zero means yes.
* `locker` - Asynchronous task process ID when the instance is in an asynchronous task.
* `machine` - Machine Model.
* `paymode` - Payment Mode.
* `period_end_time` - Instance expiration time, the format is 2006-01-02 15:04:05.
* `pid` - Product Type ID.
* `qps` - Maximum Qps value.
* `region` - The name of the region where the instance is located, such as ap-shanghai.
* `status_desc` - Description of the current running state of the instance.
* `status` - Instance status: 0 creating, 1 process processing, 2 running, 3 instance not initialized, -1 instance isolated, 4 instance initializing, 5 instance deleting, 6 instance restarting, 7 data migration.
* `tdsql_version` - TDSQL version information.
* `uin` - The account to which the instance belongs.
* `update_time` - The last update time of the instance in the format of 2006-01-02 15:04:05.
* `vipv6` - Intranet IPv6.
* `vport` - Intranet port.
* `wan_domain` - The domain name accessed from the external network, which can be resolved by the public network.
* `wan_port_ipv6` - Internet IPv6 port.
* `wan_port` - Internet port.
* `wan_status_ipv6` - Internet IPv6 status.
* `wan_status` - External network status, 0-unopened; 1-opened; 2-closed; 3-opening.
* `wan_vip` - Extranet IP address, accessible from the public network.
* `wan_vipv6` - Internet IPv6.


## Import

mariadb tencentcloud_mariadb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_instance.instance tdsql-4pzs5b67
```

