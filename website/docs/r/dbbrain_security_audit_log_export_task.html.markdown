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
resource "tencentcloud_dbbrain_security_audit_log_export_task" "task" {
  sec_audit_group_id = "sec_audit_group_id"
  start_time         = "2020-12-28 00:00:00"
  end_time           = "2020-12-28 01:00:00"
  product            = "mysql"
  danger_levels      = [0, 1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) end time.
* `product` - (Required, String, ForceNew) product, optional value is mysql.
* `sec_audit_group_id` - (Required, String, ForceNew) security audit group id.
* `start_time` - (Required, String, ForceNew) start time.
* `danger_levels` - (Optional, Set: [`Int`], ForceNew) List of log risk levels, supported values include: 0 no risk; 1 low risk; 2 medium risk; 3 high risk.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



