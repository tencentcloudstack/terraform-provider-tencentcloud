---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_modify_diag_db_instance_operation"
sidebar_current: "docs-tencentcloud-resource-dbbrain_modify_diag_db_instance_operation"
description: |-
  Provides a resource to create a dbbrain modify_diag_db_instance_conf
---

# tencentcloud_dbbrain_modify_diag_db_instance_operation

Provides a resource to create a dbbrain modify_diag_db_instance_conf

## Example Usage

```hcl
resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "on" {
  instance_confs {
    daily_inspection = "Yes"
    overview_display = "Yes"
  }
  product      = "mysql"
  instance_ids = ["%s"]
}
```
```hcl
resource "tencentcloud_dbbrain_modify_diag_db_instance_operation" "off" {
  instance_confs {
    daily_inspection = "No"
    overview_display = "No"
  }
  product      = "mysql"
  instance_ids = ["%s"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_confs` - (Required, List, ForceNew) Instance configuration, including inspection, overview switch, etc.
* `product` - (Required, String, ForceNew) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL.
* `instance_ids` - (Optional, Set: [`String`], ForceNew) Specifies the ID of the instance whose inspection status is changed.
* `regions` - (Optional, String, ForceNew) Effective instance region, the value is All, which means all regions.

The `instance_confs` object supports the following:

* `daily_inspection` - (Optional, String) Database inspection switch, Yes/No.
* `overview_display` - (Optional, String) Instance overview switch, Yes/No.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



