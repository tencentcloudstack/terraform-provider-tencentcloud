---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_rotate_product_secret"
sidebar_current: "docs-tencentcloud-resource-ssm_rotate_product_secret"
description: |-
  Provides a resource to create a ssm rotate_product_secret
---

# tencentcloud_ssm_rotate_product_secret

Provides a resource to create a ssm rotate_product_secret

## Example Usage

```hcl
resource "tencentcloud_ssm_rotate_product_secret" "example" {
  secret_name = "tf_example"
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String, ForceNew) Secret name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



