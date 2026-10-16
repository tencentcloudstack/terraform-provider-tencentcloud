---
subcategory: "Billing"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_billing_allocation_tag"
sidebar_current: "docs-tencentcloud-resource-billing_allocation_tag"
description: |-
  Provides a resource to create a Billing allocation tag
---

# tencentcloud_billing_allocation_tag

Provides a resource to create a Billing allocation tag

## Example Usage

```hcl
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}

resource "tencentcloud_billing_allocation_tag" "example" {
  tag_key = tencentcloud_tag.example.tag_key
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, String, ForceNew) Cost allocation tag key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Tag type, 0 normal tag, 1 account tag.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAllocationTag](https://www.tencentcloud.com/document/api/555/57415)
* Read: [DescribeTagList](https://www.tencentcloud.com/document/api/555/57412)
* Delete: [DeleteAllocationTag](https://www.tencentcloud.com/document/api/555/57414)

## Import

Billing allocation tag can be imported using the id, e.g.

```
terraform import tencentcloud_billing_allocation_tag.example tagKey
```

