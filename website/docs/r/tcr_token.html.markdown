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

### Create a token for tcr instance

```hcl
resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_token" "example" {
  instance_id = tencentcloud_tcr_instance.example.id
  description = "example for the tcr token"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the TCR instance.
* `description` - (Optional, String, ForceNew) Description of the token. Valid length is [0~255].
* `enable` - (Optional, Bool) Indicate to enable this token or not.

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
$ terraform import tencentcloud_tcr_token.example instance_id#token_id
```

