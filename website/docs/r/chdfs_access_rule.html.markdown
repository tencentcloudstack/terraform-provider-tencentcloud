---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_access_rule"
sidebar_current: "docs-tencentcloud-resource-chdfs_access_rule"
description: |-
  Provides a resource to create a chdfs access_rule
---

# tencentcloud_chdfs_access_rule

Provides a resource to create a chdfs access_rule

## Example Usage

```hcl
resource "tencentcloud_chdfs_access_rule" "access_rule" {
  access_group_id = "ag-bvmzrbsm"

  access_rule {
    access_mode = 2
    address     = "10.0.1.1"
    priority    = 12
  }
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Required, String, ForceNew) access group id.
* `access_rule` - (Required, List) rule detail.

The `access_rule` object supports the following:

* `access_mode` - (Optional, Int) rule access mode, 1: read only, 2: read &amp; wirte.
* `address` - (Optional, String) rule address, IP OR IP SEG.
* `priority` - (Optional, Int) rule priority, range 1 - 100, value less higher priority.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

chdfs access_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_access_rule.access_rule access_group_id#access_rule_id
```

