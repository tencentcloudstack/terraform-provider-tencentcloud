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
  cluster_id            = "cynosdbmysql-bws8h88b"
  instance_id           = "cynosdbmysql-ins-rikr6z4o"
  is_in_maintain_period = "no"

  instance_param_list {
    current_value = "0"
    param_name    = "init_connect"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `instance_id` - (Optional, String, ForceNew) Instance ID.
* `instance_param_list` - (Optional, Set) Instance parameter list.
* `is_in_maintain_period` - (Optional, String) Yes: modify within the operation and maintenance time window, no: execute immediately (default value).

The `instance_param_list` object supports the following:

* `current_value` - (Required, String) Current value of parameter.
* `param_name` - (Required, String) Parameter Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



