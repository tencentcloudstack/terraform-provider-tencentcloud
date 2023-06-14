---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_diag_db_instances"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_diag_db_instances"
description: |-
  Use this data source to query detailed information of dbbrain diag_db_instances
---

# tencentcloud_dbbrain_diag_db_instances

Use this data source to query detailed information of dbbrain diag_db_instances

## Example Usage

```hcl
data "tencentcloud_dbbrain_diag_db_instances" "diag_db_instances" {
  is_supported   = true
  product        = "mysql"
  instance_names = ["keep_preset_mysql"]
}
```

## Argument Reference

The following arguments are supported:

* `is_supported` - (Required, Bool) whether it is an instance supported by DBbrain, always pass `true`.
* `product` - (Required, String) service product type, supported values include: `mysql` - cloud database MySQL, `cynosdb` - cloud database TDSQL-C for MySQL, the default is `mysql`.
* `instance_ids` - (Optional, Set: [`String`]) query based on the instance ID condition.
* `instance_names` - (Optional, Set: [`String`]) query based on the instance name condition.
* `regions` - (Optional, Set: [`String`]) query based on geographical conditions.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_scan_status` - all-instance inspection status. `0`: All-instance inspection is enabled; `1`: All-instance inspection is not enabled.
* `items` - information about the instance.
  * `audit_policy_status` - Instance audit log enable status. `ALL_AUDIT`: full audit is enabled; `RULE_AUDIT`: rule audit is enabled; `UNBOUND`: audit is not enabled.
  * `audit_running_status` - Instance audit log running status. `normal`: running; `paused`: arrears suspended.
  * `cpu` - number of cores.
  * `create_time` - create time.
  * `deadline_time` - resource expiration time.
  * `deploy_mode` - cdb type.
  * `engine_version` - database version.
  * `event_count` - the number of abnormal events.
  * `group_id` - group ID.
  * `group_name` - group name.
  * `health_score` - health score.
  * `init_flag` - cdb instance initialization flag: `0`: not initialized; `1`: initialized.
  * `instance_conf` - status of instance inspection/overview.
    * `daily_inspection` - database inspection switch, Yes/No.
    * `key_delimiters` - Custom separator for redis large key analysis, only used by `redis`. Note: This field may return null, indicating that no valid value can be obtained.
    * `overview_display` - instance overview switch, Yes/No.
  * `instance_id` - instance id.
  * `instance_name` - instance name.
  * `instance_type` - instance type. `1`: MASTER; `2`: DR, `3`: RO, `4`: SDR.
  * `internal_vip` - Intranet VIPNote: This field may return null, indicating that no valid value can be obtained.
  * `internal_vport` - Intranet portNote: This field may return null, indicating that no valid value can be obtained.
  * `is_supported` - whether it is an instance supported by DBbrain.
  * `memory` - memory, in MB.
  * `product` - belongs to the product.
  * `region` - region.
  * `sec_audit_status` - enabled status of the instance security audit log. `ON`: security audit is enabled; `OFF`: security audit is not enabled.
  * `source` - access source.
  * `status` - Instance status: `0`: Shipping; `1`: Running normally; `4`: Destroying; `5`: Isolating.
  * `task_status` - task status.
  * `uniq_subnet_id` - subnet uniform ID.
  * `uniq_vpc_id` - the unified ID of the private network.
  * `vip` - intranet address.
  * `volume` - hard disk storage, in GB.
  * `vport` - intranet port.


