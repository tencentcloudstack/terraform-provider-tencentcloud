---
subcategory: "TencentDB for MariaDB(MariaDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mariadb_dcn_detail"
sidebar_current: "docs-tencentcloud-datasource-mariadb_dcn_detail"
description: |-
  Use this data source to query detailed information of mariadb dcn_detail
---

# tencentcloud_mariadb_dcn_detail

Use this data source to query detailed information of mariadb dcn_detail

## Example Usage

```hcl
data "tencentcloud_mariadb_dcn_detail" "dcn_detail" {
  instance_id = "tdsql-9vqvls95"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dcn_details` - DCN synchronization details.
  * `cpu` - Number of CPU cores of the instance.
  * `create_time` - Creation time of the instance in the format of 2006-01-02 15:04:05.
  * `dcn_flag` - DCN flag. Valid values: `1` (primary), `2` (disaster recovery).
  * `dcn_status` - DCN status. Valid values: `0` (none), `1` (creating), `2` (syncing), `3` (disconnected).
  * `encrypt_status` - Whether KMS is enabled.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_type` - Instance type. Valid values: `1` (dedicated primary instance), `2` (non-dedicated primary instance), `3` (non-dedicated disaster recovery instance), `4` (dedicated disaster recovery instance).
  * `memory` - Instance memory capacity in GB.
  * `pay_mode` - Billing mode.
  * `period_end_time` - Expiration time of the instance in the format of 2006-01-02 15:04:05.
  * `region` - Region where the instance resides.
  * `replica_config` - Configuration information of DCN replication. This field is null for a primary instance.Note: This field may return null, indicating that no valid values can be obtained.
    * `delay_replication_type` - Delayed replication type. Valid values: `DEFAULT` (no delay), `DUE_TIME` (specified replication time)Note: This field may return null, indicating that no valid values can be obtained.
    * `due_time` - Specified time for delayed replicationNote: This field may return null, indicating that no valid values can be obtained.
    * `replication_delay` - The number of seconds to delay the replicationNote: This field may return null, indicating that no valid values can be obtained.
    * `ro_replication_mode` - DCN running status. Valid values: `START` (running), `STOP` (pause)Note: This field may return null, indicating that no valid values can be obtained.
  * `replica_status` - DCN replication status. This field is null for the primary instance.Note: This field may return null, indicating that no valid values can be obtained.
    * `delay` - The current delay, which takes the delay value of the replica instance.
    * `status` - DCN running status. Valid values: `START` (running), `STOP` (pause).Note: This field may return null, indicating that no valid values can be obtained.
  * `status_desc` - Instance status description.
  * `status` - Instance status.
  * `storage` - Instance storage capacity in GB.
  * `vip` - Instance IP address.
  * `vipv6` - Instance IPv6 address.
  * `vport` - Instance port.
  * `zone` - Availability zone where the instance resides.


