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

The `access_rule` object exports the following:

* `access_rule_id` - single rule id.
* `create_time` - rule create time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAccessRules](https://www.tencentcloud.com/search/?wd=CreateAccessRules&type=document&lang=en)
* Read: [DescribeAccessRules](https://www.tencentcloud.com/search/?wd=DescribeAccessRules&type=document&lang=en)
* Update: [ModifyAccessRules](https://www.tencentcloud.com/search/?wd=ModifyAccessRules&type=document&lang=en)
* Delete: [DeleteAccessRules](https://www.tencentcloud.com/search/?wd=DeleteAccessRules&type=document&lang=en)

## Import

chdfs access_rule can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_access_rule.access_rule access_group_id#access_rule_id
```

