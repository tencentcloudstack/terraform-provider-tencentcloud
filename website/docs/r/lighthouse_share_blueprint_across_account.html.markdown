---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_share_blueprint_across_account"
sidebar_current: "docs-tencentcloud-resource-lighthouse_share_blueprint_across_account"
description: |-
  Provides a resource to share a Lighthouse blueprint across accounts.
---

# tencentcloud_lighthouse_share_blueprint_across_account

Provides a resource to share a Lighthouse blueprint across accounts.

## Example Usage

```hcl
resource "tencentcloud_lighthouse_share_blueprint_across_account" "example" {
  blueprint_id = "lhbp-xxxxxx"
  account_ids  = ["100000000001", "100000000002"]
}
```

## Argument Reference

The following arguments are supported:

* `account_ids` - (Required, List: [`String`], ForceNew) List of account IDs receiving the shared image. Max 10 account IDs. Must be main accounts.
* `blueprint_id` - (Required, String, ForceNew) Blueprint ID. Only custom images in NORMAL state can be shared.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



