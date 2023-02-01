---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_cluster_instance_groups"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_cluster_instance_groups"
description: |-
  Use this data source to query detailed information of cynosdb cluster_instance_groups
---

# tencentcloud_cynosdb_cluster_instance_groups

Use this data source to query detailed information of cynosdb cluster_instance_groups

## Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_instance_groups" "cluster_instance_groups" {
  cluster_id = xxxxxx ;
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) The ID of cluster.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_grp_info_list` - List of instance groups.
  * `app_id` - App id.
  * `cluster_id` - The ID of cluster.
  * `created_time` - Created time.
  * `deleted_time` - Deleted time.
  * `instance_grp_id` - The ID of instance group.
  * `instance_set` - Instance groups contain instance information.
    * `app_id` - User app id.
    * `business_type` - Business type.Note: This field may return null, indicating that no valid value can be obtained.
    * `cluster_id` - The id of cluster.
    * `cluster_name` - The name of cluster.
    * `cpu` - Cpu, unit: CORE.
    * `create_time` - Create time.
    * `cynos_version` - Cynos kernel version.
    * `db_type` - Database type.
    * `db_version` - Database version.
    * `destroy_deadline_text` - Destroy deadline.
    * `destroy_time` - Instance destroy time.
    * `instance_id` - The id of instance.
    * `instance_name` - The name of instance.
    * `instance_role` - Instance role.
    * `instance_type` - Instance type.
    * `is_freeze` - Whether to freeze.Note: This field may return null, indicating that no valid value can be obtained.
    * `isolate_time` - Isolate time.
    * `max_cpu` - Serverless instance maxmum cpu.
    * `memory` - Memory, unit: GB.
    * `min_cpu` - Serverless instance minimum cpu.
    * `net_type` - Net type.
    * `pay_mode` - Pay mode.
    * `period_end_time` - Instance expiration time.
    * `physical_zone` - Physical zone.
    * `processing_task` - Task being processed.
    * `project_id` - The id of project.
    * `region` - Region.
    * `renew_flag` - Renew flag.
    * `resource_tags` - Resource tags.Note: This field may return null, indicating that no valid value can be obtained.
      * `tag_key` - The key of tag.
      * `tag_value` - The value of tag.
    * `serverless_status` - Serverless instance status, optional values:resumepause.
    * `status_desc` - Instance state Chinese description.
    * `status` - The status of instance.
    * `storage_id` - Prepaid Storage Id.Note: This field may return null, indicating that no valid value can be obtained..
    * `storage_pay_mode` - Storage payment type.
    * `storage` - Storage, unit: GB.
    * `subnet_id` - Subnet ID.
    * `tasks` - Task list.Note: This field may return null, indicating that no valid value can be obtained.
      * `object_id` - Task ID (cluster ID|instance group ID|instance ID).Note: This field may return null, indicating that no valid value can be obtained.
      * `object_type` - Object type.Note: This field may return null, indicating that no valid value can be obtained.
      * `task_id` - Task auto-increment ID.Note: This field may return null, indicating that no valid value can be obtained.
      * `task_status` - Task status.Note: This field may return null, indicating that no valid value can be obtained.
      * `task_type` - Task type.Note: This field may return null, indicating that no valid value can be obtained.
    * `uin` - User Uin.
    * `update_time` - Update time.
    * `vip` - Instance intranet IP.
    * `vpc_id` - VPC network ID.
    * `vport` - Instance intranet VPort.
    * `wan_domain` - Public domain.
    * `wan_ip` - Public IP.
    * `wan_port` - Public port.
    * `wan_status` - Public status.
    * `zone` - Availability zone.
  * `status` - Status.
  * `type` - Instance group type. ha-ha group; ro-read-only group.
  * `updated_time` - Updated time.
  * `vip` - Intranet IP.
  * `vport` - Intranet port.
  * `wan_domain` - Public domain name.
  * `wan_ip` - Public IP.
  * `wan_port` - Public port.
  * `wan_status` - Public status.


