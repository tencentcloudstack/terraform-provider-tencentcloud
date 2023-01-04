---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_hot_link"
sidebar_current: "docs-tencentcloud-resource-ci_hot_link"
description: |-
  Provides a resource to create a ci hot_link
---

# tencentcloud_ci_hot_link

Provides a resource to create a ci hot_link

## Example Usage

```hcl
resource "tencentcloud_ci_hot_link" "hot_link" {
  bucket = "terraform-ci-xxxxxx"
  url    = ["10.0.0.1", "10.0.0.2"]
  type   = "white"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `type` - (Required, String) Anti-leech type, `white` is whitelist, `black` is blacklist.
* `url` - (Required, Set: [`String`]) domain address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci hot_link can be imported using the bucket, e.g.

```
terraform import tencentcloud_ci_hot_link.hot_link terraform-ci-xxxxxx
```

