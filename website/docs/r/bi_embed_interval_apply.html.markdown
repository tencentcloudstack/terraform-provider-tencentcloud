---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_embed_interval_apply"
sidebar_current: "docs-tencentcloud-resource-bi_embed_interval_apply"
description: |-
  Provides a resource to create a bi embed_interval
---

# tencentcloud_bi_embed_interval_apply

Provides a resource to create a bi embed_interval

## Example Usage

```hcl
resource "tencentcloud_bi_embed_interval_apply" "embed_interval" {
  project_id = 11015030
  page_id    = 10520483
  bi_token   = "4192d65b-d674-4117-9a59-xxxxxxxxx"
  scope      = "page"
}
```

## Argument Reference

The following arguments are supported:

* `bi_token` - (Optional, String, ForceNew) Token that needs to be applied for extension.
* `page_id` - (Optional, Int, ForceNew) Sharing page id, this is empty value 0 when embedding the board.
* `project_id` - (Optional, Int, ForceNew) Sharing project id, required.
* `scope` - (Optional, String, ForceNew) Choose panel or page.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



