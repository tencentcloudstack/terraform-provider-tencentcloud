---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_group"
sidebar_current: "docs-tencentcloud-resource-security_group"
description: |-
  Provides a resource to create security group.
---

# tencentcloud_security_group

Provides a resource to create security group.

## Example Usage

### Create a basic security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example-sg"
  description = "sg test"
}
```

### Create a complete security group

```hcl
resource "tencentcloud_security_group" "example" {
  name        = "tf-example-sg"
  description = "sg test"
  project_id  = 0

  tags = {
    "example" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the security group to be queried.
* `description` - (Optional, String) Description of the security group.
* `project_id` - (Optional, Int, ForceNew) Project ID of the security group.
* `tags` - (Optional, Map) Tags of the security group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Security group can be imported using the id, e.g.

```
  $ terraform import tencentcloud_security_group.sglab sg-ey3wmiz1
```

