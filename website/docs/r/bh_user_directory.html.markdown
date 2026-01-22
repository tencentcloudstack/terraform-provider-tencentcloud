---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_user_directory"
sidebar_current: "docs-tencentcloud-resource-bh_user_directory"
description: |-
  Provides a resource to create a BH user directory
---

# tencentcloud_bh_user_directory

Provides a resource to create a BH user directory

## Example Usage

```hcl
resource "tencentcloud_bh_user_directory" "example" {
  dir_id   = 895784
  dir_name = "tf-example"
  user_org_set {
    org_id        = 1576799
    org_name      = "orgName1"
    org_id_path   = "819729.895784"
    org_name_path = "Root.demo1"
    user_total    = 0
  }

  user_org_set {
    org_id        = 896536
    org_name      = "orgName2"
    org_id_path   = "819729.895784.896536"
    org_name_path = "Root.demo2.demo3"
    user_total    = 1
  }
  source      = 0
  source_name = "sourceName"
}
```

## Argument Reference

The following arguments are supported:

* `dir_id` - (Required, Int, ForceNew) Directory ID.
* `dir_name` - (Required, String, ForceNew) Directory name.
* `source_name` - (Required, String, ForceNew) IOA associated user source name.
* `source` - (Required, Int, ForceNew) IOA associated user source type.
* `user_org_set` - (Required, List) IOA group information.

The `user_org_set` object supports the following:

* `org_id_path` - (Required, String) IOA user organization ID path.
* `org_id` - (Required, Int) IOA user organization ID.
* `org_name_path` - (Required, String) IOA user organization name path.
* `org_name` - (Required, String) IOA user organization name.
* `user_total` - (Optional, Int) Number of users under the IOA user organization ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `directory_id` - Directory ID.
* `user_count` - Number of users included in the directory.


## Import

BH user directory can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user_directory.example 32
```

