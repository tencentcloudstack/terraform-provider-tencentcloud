---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_switch_for_upgrade"
sidebar_current: "docs-tencentcloud-resource-mysql_switch_for_upgrade"
description: |-
  Provides a resource to create a mysql switch_for_upgrade
---

# tencentcloud_mysql_switch_for_upgrade

Provides a resource to create a mysql switch_for_upgrade

## Example Usage

```hcl
resource "tencentcloud_mysql_switch_for_upgrade" "switch_for_upgrade" {
  instance_id = "cdb-d9gbh7lt"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed on the TencentDB Console page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



