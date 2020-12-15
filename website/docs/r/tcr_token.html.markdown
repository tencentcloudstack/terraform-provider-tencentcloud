---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_token"
sidebar_current: "docs-tencentcloud-resource-tcr_token"
description: |-
  Use this resource to create tcr long term token.
---

# tencentcloud_tcr_token

Use this resource to create tcr long term token.

## Example Usage

```hcl
resource "tencentcloud_tcr_token" "foo" {
  instance_id = "cls-cda1iex1"
  description = "test"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the TCR instance.
* `description` - (Optional, ForceNew) Description of the token. Valid length is [0~255].
* `enable` - (Optional) Indicate to enable this token or not.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `token_id` - Sub ID of the TCR token. The full ID of token format like `instance_id#token_id`.
* `token` - The content of the token.
* `user_name` - User name of the token.


## Import

tcr token can be imported using the id, e.g.

```
$ terraform import tencentcloud_tcr_token.foo cls-cda1iex1#namespace#buv3h3j96j2d1rk1cllg
```

