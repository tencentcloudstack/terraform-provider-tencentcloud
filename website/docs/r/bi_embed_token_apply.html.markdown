---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_embed_token_apply"
sidebar_current: "docs-tencentcloud-resource-bi_embed_token_apply"
description: |-
  Provides a resource to create a bi embed_token
---

# tencentcloud_bi_embed_token_apply

Provides a resource to create a bi embed_token

## Example Usage

```hcl
resource "tencentcloud_bi_embed_token_apply" "embed_token" {
  project_id   = 11015030
  page_id      = 10520483
  scope        = "page"
  expire_time  = "240"
  user_corp_id = "100022975249"
  user_id      = "100024664626"
}
```

## Argument Reference

The following arguments are supported:

* `expire_time` - (Optional, String, ForceNew) Expiration. Unit: Minutes Maximum value: 240. i.e. 4 hours Default: 240.
* `page_id` - (Optional, Int, ForceNew) Sharing page id, this is empty value 0 when embedding the board.
* `project_id` - (Optional, Int, ForceNew) Share project id.
* `scope` - (Optional, String, ForceNew) Page means embedding the page, and panel means embedding the entire board.
* `user_corp_id` - (Optional, String, ForceNew) User enterprise ID (for multi-user only).
* `user_id` - (Optional, String, ForceNew) UserId (for multi-user only).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bi_token` - Create the generated token.
* `create_at` - Create time.
* `udpate_at` - Upadte time.


