---
subcategory: "PostgreSQL"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_readonly_attachment"
sidebar_current: "docs-tencentcloud-resource-postgresql_readonly_attachment"
description: |-
  Use this resource to create postgresql readonly attachment.
---

# tencentcloud_postgresql_readonly_attachment

Use this resource to create postgresql readonly attachment.

## Example Usage

```hcl
resource "tencentcloud_postgresql_readonly_attachment" "attach" {
  db_instance_id     = tencentcloud_postgresql_readonly_instance.foo.id
  read_only_group_id = tencentcloud_postgresql_readonly_group.group.id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Read only instance ID.
* `read_only_group_id` - (Required, String, ForceNew) Read only group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



