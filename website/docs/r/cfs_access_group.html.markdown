---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_access_group"
sidebar_current: "docs-tencentcloud-resource-cfs_access_group"
description: |-
  Provides a resource to create a CFS access group.
---

# tencentcloud_cfs_access_group

Provides a resource to create a CFS access group.

## Example Usage

```hcl
resource "tencentcloud_cfs_access_group" "example" {
  name        = "tx_example"
  description = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the access group, and max length is 64.
* `description` - (Optional, String) Description of the access group, and max length is 255.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the access group.


## Import

CFS access group can be imported using the id, e.g.

```
$ terraform import tencentcloud_cfs_access_group.example pgroup-7nx89k7l
```

