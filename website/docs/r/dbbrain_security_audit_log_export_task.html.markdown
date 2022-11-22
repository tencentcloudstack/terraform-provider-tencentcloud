---
subcategory: "DBbrain"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_security_audit_log_export_task"
sidebar_current: "docs-tencentcloud-resource-dbbrain_security_audit_log_export_task"
description: |-
  Provides a resource to create a dbbrain security_audit_log_export_task
---

# tencentcloud_dbbrain_security_audit_log_export_task

Provides a resource to create a dbbrain security_audit_log_export_task

## Example Usage

```hcl
resource "tencentcloud_dbbrain_security_audit_log_export_task" "security_audit_log_export_task" {
  sec_audit_group_id = ""
  start_time         = ""
  end_time           = ""
  product            = ""
  danger_levels      = ""
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) end time.
* `product` - (Required, String) product, optional value is mysql.
* `sec_audit_group_id` - (Required, String) security audit group id.
* `start_time` - (Required, String) start time.
* `danger_levels` - (Optional, Set: [`Int`]) List of log risk levels, supported values include: 0 no risk; 1 low risk; 2 medium risk; 3 high risk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dbbrain security_audit_log_export_task can be imported using the id, e.g.
```
$ terraform import tencentcloud_dbbrain_security_audit_log_export_task.security_audit_log_export_task securityAuditLogExportTask_id
```

