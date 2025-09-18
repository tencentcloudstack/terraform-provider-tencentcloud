---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_audit_service"
sidebar_current: "docs-tencentcloud-resource-mysql_audit_service"
description: |-
  Provides a resource to create a Mysql audit service
---

# tencentcloud_mysql_audit_service

Provides a resource to create a Mysql audit service

## Example Usage

### If audit_all is true

```hcl
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
```

### If audit_all is false

```hcl
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 30
  high_log_expire_day = 7
  rule_template_ids = [
    "cdb-art-3a9ww0oj"
  ]
  audit_all = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) TencentDB for MySQL instance ID.
* `log_expire_day` - (Required, Int) Retention period of the audit log. Valid values:  `7` (one week), `30` (one month), `90` (three months), `180` (six months), `365` (one year), `1095` (three years), `1825` (five years).
* `audit_all` - (Optional, Bool) Audit type. Valid values: true: Record all; false: Record by rules (default value).
* `high_log_expire_day` - (Optional, Int) Retention period of high-frequency audit logs. Valid values:  `7` (one week), `30` (one month).
* `rule_template_ids` - (Optional, Set: [`String`]) Rule template ID. If both this parameter and AuditRuleFilters are not specified, all SQL statements will be recorded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Mysql audit service can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_audit_service.example cdb-3kwa3gfj
```

