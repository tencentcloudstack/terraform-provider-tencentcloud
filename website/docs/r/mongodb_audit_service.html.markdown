---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_audit_service"
sidebar_current: "docs-tencentcloud-resource-mongodb_audit_service"
description: |-
  Provides a resource to manage MongoDB audit service.
---

# tencentcloud_mongodb_audit_service

Provides a resource to manage MongoDB audit service.

## Example Usage

### Full audit mode

```hcl
resource "tencentcloud_mongodb_audit_service" "example" {
  instance_id    = "cmgo-5aqo4yf7"
  log_expire_day = 7
  audit_all      = true
}
```

### Rule-based audit mode

```hcl
resource "tencentcloud_mongodb_audit_service" "example" {
  instance_id    = "cmgo-5aqo4yf7"
  log_expire_day = 30
  audit_all      = false

  rule_filters {
    type    = "DB"
    compare = "EQ"
    value   = ["testdb"]
  }

  rule_filters {
    type    = "User"
    compare = "EQ"
    value   = ["admin"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `audit_all` - (Required, Bool) Whether to enable full audit. true: full audit, false: rule-based audit. When AuditAll is true, RuleFilters is not required.
* `instance_id` - (Required, String, ForceNew) Instance ID, for example: cmgo-xfts****.
* `log_expire_day` - (Required, Int) Audit log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.
* `rule_filters` - (Optional, List) Audit filter rules. Only required when audit_all is false.

The `rule_filters` object supports the following:

* `compare` - (Required, String) Filter match type. Must be EQ.
* `type` - (Required, String) Filter condition name. Valid values: SrcIp, DB, Collection, User, SqlType.
* `value` - (Required, List) Filter match values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_name` - Instance name.
* `log_type` - Audit log storage type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `5m`) Used when creating the resource.
* `delete` - (Defaults to `5m`) Used when deleting the resource.

## Import

MongoDB audit service can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_audit_service.example cmgo-5aqo4yf7
```

