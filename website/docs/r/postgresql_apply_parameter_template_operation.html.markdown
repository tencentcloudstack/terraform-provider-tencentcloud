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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeDBInstanceAttribute](https://www.tencentcloud.com/document/product/409/16772) - This API is used to query the details of one instance.
* Create: [DescribeDBInstanceParameters](https://www.tencentcloud.com/document/product/409/43284) - This API is used to query the parameters of an instance.
* Create: [DescribeParameterTemplateAttributes](https://www.tencentcloud.com/document/product/409/52652) - This API is used to query the details of a parameter template, including basic information and parameter information.
* Create: [DescribeParamsEvent](https://www.tencentcloud.com/document/product/409/43283) - This API is used to query the parameter modification event.
* Create: [ModifyDBInstanceParameters](https://www.tencentcloud.com/document/product/409/43282) - This API is used to modify instance parameters.

