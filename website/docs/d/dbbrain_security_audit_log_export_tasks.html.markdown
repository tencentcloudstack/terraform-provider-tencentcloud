---
subcategory: "DBbrain"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_security_audit_log_export_tasks"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_security_audit_log_export_tasks"
description: |-
  Use this data source to query detailed information of dbbrain securityAuditLogExportTasks
---

# tencentcloud_dbbrain_security_audit_log_export_tasks

Use this data source to query detailed information of dbbrain securityAuditLogExportTasks

## Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "sec_audit_group_id"
  start_time         = "start_time"
  end_time           = "end_time"
  product            = "mysql"
  danger_levels      = [0, 1, 2]
}

data "tencentcloud_dbbrain_security_audit_log_export_tasks" "tasks" {
  sec_audit_group_id = "sec_audit_group_id"
  product            = "mysql"
  async_request_ids  = [tencentcloud_dbbrain_security_audit_log_export_task.task.async_request_id]
}
```

## Argument Reference

The following arguments are supported:

* `product` - (Required, String) product, optional value is mysql.
* `sec_audit_group_id` - (Required, String) security audit group id.
* `async_request_ids` - (Optional, Set: [`Int`]) async request id list.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - security audit log export task list.
  * `async_request_id` - async request id.
  * `create_time` - create time.
  * `danger_levels` - danger level list.
  * `end_time` - end time.
  * `log_end_time` - log end time.
  * `log_start_time` - log start time.
  * `progress` - task progress.
  * `start_time` - start time.
  * `status` - status.
  * `total_size` - the total size of log.


