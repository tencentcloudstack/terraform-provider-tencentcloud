---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_audit_service"
sidebar_current: "docs-tencentcloud-resource-cynosdb_audit_service"
description: |-
  Provides a resource to create a CynosDB audit service
---

# tencentcloud_cynosdb_audit_service

Provides a resource to create a CynosDB audit service

## Example Usage

### If audit_all is true

```hcl
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-ins-f9j6sopi"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
```

### If audit_all is false

```hcl
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-ins-f9j6sopi"
  log_expire_day      = 30
  high_log_expire_day = 7
  rule_template_ids   = ["cynosdb-art-riwq2vx0"]
  audit_all           = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `log_expire_day` - (Required, Int) Log retention period.
* `audit_all` - (Optional, Bool) Audit type. true - full audit; default false - rule-based audit.
* `high_log_expire_day` - (Optional, Int) Frequent log retention period.
* `rule_template_ids` - (Optional, Set: [`String`]) Rule template ID set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CynosDB audit service can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_audit_service.example cynosdbmysql-ins-f9j6sopi
```

