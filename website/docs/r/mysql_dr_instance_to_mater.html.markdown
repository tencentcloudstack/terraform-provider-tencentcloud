---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_dr_instance_to_mater"
sidebar_current: "docs-tencentcloud-resource-mysql_dr_instance_to_mater"
description: |-
  Provides a resource to create a mysql dr_instance_to_mater
---

# tencentcloud_mysql_dr_instance_to_mater

Provides a resource to create a mysql dr_instance_to_mater

## Example Usage

```hcl
resource "tencentcloud_mysql_dr_instance_to_mater" "dr_instance_to_mater" {
  instance_id = "cdb-c1nl9rpv"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Disaster recovery instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed in the TencentDB console.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql dr_instance_to_mater can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_dr_instance_to_mater.dr_instance_to_mater dr_instance_to_mater_id
```

