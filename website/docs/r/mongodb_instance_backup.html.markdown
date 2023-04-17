---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_backup"
sidebar_current: "docs-tencentcloud-resource-mongodb_instance_backup"
description: |-
  Provides a resource to create a mongodb instance_backup
---

# tencentcloud_mongodb_instance_backup

Provides a resource to create a mongodb instance_backup

## Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup" "instance_backup" {
  instance_id   = "cmgo-9d0p6umb"
  backup_method = 0
  backup_remark = "my backup"
}
```

## Argument Reference

The following arguments are supported:

* `backup_method` - (Required, Int, ForceNew) 0:logical backup, 1:physical backup.
* `instance_id` - (Required, String, ForceNew) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `backup_remark` - (Optional, String, ForceNew) backup notes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



