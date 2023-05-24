---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_instance_param"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_instance_param"
description: |-
  Provides a resource to create a sqlserver config_instance_param
---

# tencentcloud_sqlserver_config_instance_param

Provides a resource to create a sqlserver config_instance_param

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_instance_param" "config_instance_param" {
  instance_id = tencentcloud_sqlserver_instance.test.id
  param_list {
    name          = "fill factor(%)"
    current_value = "90"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `param_list` - (Required, List) List of modified parameters. Each list element has two fields: Name and CurrentValue. Set Name to the parameter name and CurrentValue to the new value after modification. Note: if the instance needs to be restarted for the modified parameter to take effect, it will be restarted immediately or during the maintenance time. Before you modify a parameter, you can use the DescribeInstanceParams API to query whether the instance needs to be restarted.

The `param_list` object supports the following:

* `current_value` - (Optional, String) Parameter value.
* `name` - (Optional, String) Parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_instance_param can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_instance_param.config_instance_param config_instance_param_id
```

