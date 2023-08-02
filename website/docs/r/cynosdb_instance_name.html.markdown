---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_instance_name"
sidebar_current: "docs-tencentcloud-resource-cynosdb_instance_name"
description: |-
  Provides a resource to create a cynosdb instance_name
---

# tencentcloud_cynosdb_instance_name

Provides a resource to create a cynosdb instance_name

## Example Usage

```hcl
resource "tencentcloud_cynosdb_instance_name" "instance_name" {
  instance_id   = "cynosdb-ins-dokydbam"
  instance_name = "newName"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `instance_name` - (Required, String) Instance Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cynosdb instance_name can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_instance_name.instance_name instance_name_id
```

