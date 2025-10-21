---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_groups"
sidebar_current: "docs-tencentcloud-datasource-postgresql_readonly_groups"
description: |-
  Use this data source to query detailed information of postgresql read_only_groups
---

# tencentcloud_postgresql_readonly_groups

Use this data source to query detailed information of postgresql read_only_groups

## Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_group" "group" {
  master_db_instance_id       = "postgres-gzg9jb2n"
  name                        = "test-datasource"
  project_id                  = 0
  vpc_id                      = "vpc-86v957zb"
  subnet_id                   = "subnet-enm92y0m"
  replay_lag_eliminate        = 1
  replay_latency_eliminate    = 1
  max_replay_lag              = 100
  max_replay_latency          = 512
  min_delay_eliminate_reserve = 1
}

data "tencentcloud_postgresql_readonly_groups" "read_only_groups" {
  filters {
    name   = "db-master-instance-id"
    values = [tencentcloud_postgresql_readonly_group.group.master_db_instance_id]
  }
  order_by      = "CreateTime"
  order_by_type = "asc"
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter condition. The primary ID must be specified in the format of db-master-instance-id to filter results, or else null will be returned.
* `order_by_type` - (Optional, String) Sorting order. Valid values:desc, asc.
* `order_by` - (Optional, String) Sorting criterion. Valid values:ROGroupId, CreateTime, Name.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Optional, String) Filter name.
* `values` - (Optional, Set) One or more filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `read_only_group_list` - list of read-only groups.
  * `db_instance_net_info` - network information.
    * `address` - DNS domain name.
    * `ip` - IP address.
    * `net_type` - network type, 1. inner (intranet address of the basic network); 2. private (intranet address of the private network); 3. public (extranet address of the basic network or private network);.
    * `port` - connection port address.
    * `protocol_type` - The protocol type for connecting to the database, currently supported: postgresql, mssql (MSSQL compatible syntax)Note: This field may return null, indicating that no valid value can be obtained.
    * `status` - network connection status, 1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);.
    * `subnet_id` - subnet IDNote: This field may return null, indicating that no valid value can be obtained.
    * `vpc_id` - private network IDNote: This field may return null, indicating that no valid value can be obtained.
  * `master_db_instance_id` - master instance idNote: This field may return null, indicating that no valid value can be obtained.
  * `max_replay_lag` - delay time size threshold.
  * `max_replay_latency` - delay space size threshold.
  * `min_delay_eliminate_reserve` - Minimum Number of Reserved InstancesNote: This field may return null, indicating that no valid value can be obtained.
  * `network_access_list` - Read-only list of group network information (this field is obsolete)Note: This field may return null, indicating that no valid value can be obtained.
    * `resource_id` - Network resource id, instance id or RO group idNote: This field may return null, indicating that no valid value can be obtained.
    * `resource_type` - Resource type, 1-instance 2-RO groupNote: This field may return null, indicating that no valid value can be obtained.
    * `subnet_id` - subnet IDNote: This field may return null, indicating that no valid value can be obtained.
    * `vip6` - IPV6 addressNote: This field may return null, indicating that no valid value can be obtained.
    * `vip` - IPV4 addressNote: This field may return null, indicating that no valid value can be obtained.
    * `vpc_id` - private network IDNote: This field may return null, indicating that no valid value can be obtained.
    * `vpc_status` - Network status, 1-applying, 2-using, 3-deleting, 4-deletedNote: This field may return null, indicating that no valid value can be obtained.
    * `vport` - access portNote: This field may return null, indicating that no valid value can be obtained.
  * `project_id` - project idNote: This field may return null, indicating that no valid value can be obtained.
  * `read_only_db_instance_list` - instance details.
    * `app_id` - user&#39;s AppId.
    * `auto_renew` - auto-renew, 1: auto-renew, 0: no auto-renew.
    * `create_time` - instance creation time.
    * `db_charset` - instance DB character set.
    * `db_engine_config` - Configuration information for the database engineNote: This field may return null, indicating that no valid value can be obtained.
    * `db_engine` - Database engine that supports:1. postgresql (cloud database PostgreSQL);2. mssql_compatible (MSSQL compatible - cloud database PostgreSQL);Note: This field may return null, indicating that no valid value can be obtained.
    * `db_instance_class` - sales specification ID.
    * `db_instance_cpu` - the number of CPUs allocated by the instance.
    * `db_instance_id` - instance ID.
    * `db_instance_memory` - the memory size allocated by the instance, unit: GB.
    * `db_instance_name` - instance name.
    * `db_instance_net_info` - instance network connection information.
      * `address` - DNS domain name.
      * `ip` - IP address.
      * `net_type` - network type, 1. inner (intranet address of the basic network); 2. private (intranet address of the private network); 3. public (extranet address of the basic network or private network);.
      * `port` - connection port address.
      * `protocol_type` - The protocol type for connecting to the database, currently supported: postgresql, mssql (MSSQL compatible syntax)Note: This field may return null, indicating that no valid value can be obtained.
      * `status` - network connection status, 1. initing (unopened); 2. opened (opened); 3. closed (closed); 4. opening (opening); 5. closing (closed);.
      * `subnet_id` - subnet IDNote: This field may return null, indicating that no valid value can be obtained.
      * `vpc_id` - private network IDNote: This field may return null, indicating that no valid value can be obtained.
    * `db_instance_status` - Instance status, respectively: applying (applying), init (to be initialized), initing (initializing), running (running), limited run (limited run), isolated (isolated), recycling (recycling ), recycled (recycled), job running (task execution), offline (offline), migrating (migration), expanding (expanding), waitSwitch (waiting for switching), switching (switching), readonly (read-only ), restarting (restarting), network changing (network changing), upgrading (kernel version upgrade).
    * `db_instance_storage` - the size of the storage space allocated by the instance, unit: GB.
    * `db_instance_type` - instance type, the types are: 1. primary (primary instance); 2. readonly (read-only instance); 3. guard (disaster recovery instance); 4. temp (temporary instance).
    * `db_instance_version` - instance version, currently only supports standard (dual machine high availability version, one master and one slave).
    * `db_kernel_version` - Database kernel versionNote: This field may return null, indicating that no valid value can be obtained.
    * `db_major_version` - PostgreSQL major versionNote: This field may return null, indicating that no valid value can be obtained.
    * `db_node_set` - Instance node informationNote: This field may return null, indicating that no valid value can be obtained.
      * `role` - Node type, the value can be:Primary, representing the primary node;Standby, stands for standby node.
      * `zone` - Availability zone where the node is located, such as ap-guangzhou-1.
    * `db_version` - PostgreSQL version.
    * `expire_time` - instance expiration time.
    * `is_support_t_d_e` - Whether the instance supports TDE data encryption 0: not supported, 1: supportedNote: This field may return null, indicating that no valid value can be obtained.
    * `isolated_time` - instance isolation time.
    * `master_db_instance_id` - Master instance information, only returned when the instance is read-onlyNote: This field may return null, indicating that no valid value can be obtained.
    * `network_access_list` - Instance network information list (this field is obsolete)Note: This field may return null, indicating that no valid value can be obtained.
      * `resource_id` - Network resource id, instance id or RO group idNote: This field may return null, indicating that no valid value can be obtained.
      * `resource_type` - Resource type, 1-instance 2-RO groupNote: This field may return null, indicating that no valid value can be obtained.
      * `subnet_id` - subnet IDNote: This field may return null, indicating that no valid value can be obtained.
      * `vip6` - IPV6 addressNote: This field may return null, indicating that no valid value can be obtained.
      * `vip` - IPV4 addressNote: This field may return null, indicating that no valid value can be obtained.
      * `vpc_id` - private network IDNote: This field may return null, indicating that no valid value can be obtained.
      * `vpc_status` - Network status, 1-applying, 2-using, 3-deleting, 4-deletedNote: This field may return null, indicating that no valid value can be obtained.
      * `vport` - access portNote: This field may return null, indicating that no valid value can be obtained.
    * `offline_time` - offline timeNote: This field may return null, indicating that no valid value can be obtained.
    * `pay_type` - billing mode, 1. prepaid (subscription, prepaid); 2. postpaid (billing by volume, postpaid).
    * `project_id` - project ID.
    * `read_only_instance_num` - Number of read-only instancesNote: This field may return null, indicating that no valid value can be obtained.
    * `region` - The region to which the instance belongs, such as: ap-guangzhou, corresponding to the Region field of the RegionSet.
    * `status_in_readonly_group` - Status of the read-only instance in the read-only groupNote: This field may return null, indicating that no valid value can be obtained.
    * `subnet_id` - subnet ID.
    * `support_ipv6` - Whether the instance supports Ipv6, 1: support, 0: not support.
    * `tag_list` - Label information bound to the instanceNote: This field may return null, indicating that no valid value can be obtained.
      * `tag_key` - label key.
      * `tag_value` - tag value.
    * `type` - machine type.
    * `uid` - Uid of the instance.
    * `update_time` - The time when the instance performed the last update.
    * `vpc_id` - private network ID.
    * `zone` - Availability zone to which the instance belongs, such as: ap-guangzhou-3, corresponding to the Zone field of ZoneSet.
  * `read_only_group_id` - read-only group idNote: This field may return null, indicating that no valid value can be obtained.
  * `read_only_group_name` - read-only group nameNote: This field may return null, indicating that no valid value can be obtained.
  * `rebalance` - automatic load balancing switch.
  * `region` - region id.
  * `replay_lag_eliminate` - delay time switch.
  * `replay_latency_eliminate` - delay size switch.
  * `status` - state.
  * `subnet_id` - subnet-idNote: This field may return null, indicating that no valid value can be obtained.
  * `vpc_id` - virtual network id.
  * `zone` - region id.


