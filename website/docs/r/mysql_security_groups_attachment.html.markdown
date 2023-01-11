---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_security_groups_attachment"
sidebar_current: "docs-tencentcloud-resource-mysql_security_groups_attachment"
description: |-
  Provides a resource to create a mysql security_groups_attachment
---

# tencentcloud_mysql_security_groups_attachment

Provides a resource to create a mysql security_groups_attachment

## Example Usage

```hcl
resource "tencentcloud_mysql_security_groups_attachment" "security_groups_attachment" {
  security_group_id = "sg-baxfiao5"
  instance_id       = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The id of instance.
* `security_group_id` - (Required, String, ForceNew) The ID of security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql security_groups_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_security_groups_attachment.security_groups_attachment securityGroupId#instanceId
```

