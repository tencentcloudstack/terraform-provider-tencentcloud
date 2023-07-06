---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket_generate_inventory_immediately_operation"
sidebar_current: "docs-tencentcloud-resource-cos_bucket_generate_inventory_immediately_operation"
description: |-
  Provides a resource to generate a cos bucket inventory immediately
---

# tencentcloud_cos_bucket_generate_inventory_immediately_operation

Provides a resource to generate a cos bucket inventory immediately

## Example Usage

```hcl
resource "tencentcloud_cos_bucket_generate_inventory_immediately_operation" "generate_inventory_immediately" {
  inventory_id = "test"
  bucket       = "keep-test-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) Bucket.
* `inventory_id` - (Required, String, ForceNew) The id of inventory.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



