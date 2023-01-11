---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_deploy_group"
sidebar_current: "docs-tencentcloud-resource-mysql_deploy_group"
description: |-
  Provides a resource to create a mysql deploy_group
---

# tencentcloud_mysql_deploy_group

Provides a resource to create a mysql deploy_group

## Example Usage

```hcl
resource "tencentcloud_mysql_deploy_group" "deploy_group" {
  deploy_group_name = "terrform-deploy"
  description       = "deploy test"
  limit_num         = 1
  dev_class         = ["TS85"]
}
```

## Argument Reference

The following arguments are supported:

* `deploy_group_name` - (Required, String) The name of deploy group. the maximum length cannot exceed 60 characters.
* `description` - (Optional, String) The description of deploy group. the maximum length cannot exceed 200 characters.
* `dev_class` - (Optional, Set: [`String`]) The device class of deploy group. optional value is SH12+SH02, TS85, etc.
* `limit_num` - (Optional, Int) The limit on the number of instances on the same physical machine in deploy group affinity policy 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql deploy_group can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_deploy_group.deploy_group deploy_group_id
```

