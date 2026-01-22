---
subcategory: "Tag"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tag_keys"
sidebar_current: "docs-tencentcloud-datasource-tag_keys"
description: |-
  Use this data source to query detailed information of Tag keys
---

# tencentcloud_tag_keys

Use this data source to query detailed information of Tag keys

## Example Usage

### Qeury all tag keys

```hcl
data "tencentcloud_tag_keys" "tags" {}
```

### Qeury tag keys by filter

```hcl
data "tencentcloud_tag_keys" "tags" {
  create_uin   = "1486445011341"
  show_project = 1
  category     = "All"
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Optional, String) Tag type. Valid values: Custom: custom tag; System: system tag; All: all tags. Default value: All.
* `create_uin` - (Optional, Int) Creator `Uin`. If not specified, `Uin` is only used as the query condition.
* `result_output_file` - (Optional, String) Used to save results.
* `show_project` - (Optional, Int) Whether to show project. Allow values: 0: no, 1: yes.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - Tag list.


