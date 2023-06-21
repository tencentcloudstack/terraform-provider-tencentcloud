---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_instance_param"
sidebar_current: "docs-tencentcloud-resource-cynosdb_instance_param"
description: |-
  Provides a resource to create a cynosdb instance_param
---

# tencentcloud_cynosdb_instance_param

Provides a resource to create a cynosdb instance_param

## Example Usage

```hcl
resource "tencentcloud_cynosdb_instance_param" "instance_param" {
  cluster_id   = ""
  instance_ids =
  cluster_param_list {
    param_name    = ""
    current_value = ""
    old_value     = ""

  }
  instance_param_list {
    param_name    = ""
    current_value = ""
    old_value     = ""

  }
  is_in_maintain_period = ""
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `instance_id` - (Optional, String) Instance ID.
* `instance_param_list` - (Optional, List) Instance parameter list.
* `is_in_maintain_period` - (Optional, String) Yes: modify within the operation and maintenance time window, no: execute immediately (default value).

The `instance_param_list` object supports the following:

* `current_value` - (Required, String) Current value of parameter.
* `param_name` - (Required, String) Parameter Name.
* `old_value` - (Optional, String) Parameter old value (only useful when generating parameters) Note: This field may return null, indicating that a valid value cannot be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



