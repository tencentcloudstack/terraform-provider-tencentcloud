---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_organization_resource_to_share_member"
sidebar_current: "docs-tencentcloud-datasource-organization_resource_to_share_member"
description: |-
  Use this data source to query detailed information of Organization resource to share member
---

# tencentcloud_organization_resource_to_share_member

Use this data source to query detailed information of Organization resource to share member

## Example Usage

```hcl
data "tencentcloud_organization_resource_to_share_member" "example" {
  area                 = "ap-guangzhou"
  search_key           = "tf-example"
  type                 = "CVM"
  product_resource_ids = ["ins-69hg2ze0", "ins-0cxjwrog"]
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Area.
* `product_resource_ids` - (Optional, Set: [`String`]) Business resource ID. Maximum 50.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Search keywords, support business resource ID search.
* `type` - (Optional, String) Resource Type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Details.


