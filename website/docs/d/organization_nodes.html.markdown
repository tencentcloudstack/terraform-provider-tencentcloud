---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_nodes"
sidebar_current: "docs-tencentcloud-datasource-organization_nodes"
description: |-
  Use this data source to query detailed information of organization nodes
---

# tencentcloud_organization_nodes

Use this data source to query detailed information of organization nodes

## Example Usage

```hcl
data "tencentcloud_organization_nodes" "organization_nodes" {
  tags {
    tag_key   = "createBy"
    tag_value = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, List) Department tag search list, with a maximum of 10.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - List details.


