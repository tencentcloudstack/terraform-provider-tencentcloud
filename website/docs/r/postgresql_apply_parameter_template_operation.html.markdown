---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_apply_parameter_template_operation"
sidebar_current: "docs-tencentcloud-resource-postgresql_apply_parameter_template_operation"
description: |-
  Provides a resource to apply parameter template
---

# tencentcloud_postgresql_apply_parameter_template_operation

Provides a resource to apply parameter template

## Example Usage

```hcl
resource tencentcloud_postgresql_apply_parameter_template_operation "apply_parameter_template_operation" {
  db_instance_id = "postgres-xxxxxx"
  template_id    = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Instance ID.
* `template_id` - (Required, String, ForceNew) Template ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



