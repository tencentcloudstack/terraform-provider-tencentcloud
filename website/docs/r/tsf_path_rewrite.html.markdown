---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_path_rewrite"
sidebar_current: "docs-tencentcloud-resource-tsf_path_rewrite"
description: |-
  Provides a resource to create a tsf path_rewrite
---

# tencentcloud_tsf_path_rewrite

Provides a resource to create a tsf path_rewrite

## Example Usage

```hcl
resource "tencentcloud_tsf_path_rewrite" "path_rewrite" {
  gateway_group_id = "group-a2j9zxpv"
  regex            = "/test"
  replacement      = "/tt"
  blocked          = "N"
  order            = 2
}
```

## Argument Reference

The following arguments are supported:

* `blocked` - (Required, String) Whether to shield the mapped path, Y: Yes N: No.
* `gateway_group_id` - (Required, String) gateway deployment group ID.
* `order` - (Required, Int) rule order, the smaller the higher the priority.
* `regex` - (Required, String) regular expression.
* `replacement` - (Required, String) content to replace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `path_rewrite_id` - path rewrite rule ID.


## Import

tsf path_rewrite can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_path_rewrite.path_rewrite rewrite-nygq33v2
```

