---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_start_instance"
sidebar_current: "docs-tencentcloud-resource-lighthouse_start_instance"
description: |-
  Provides a resource to create a lighthouse start_instance
---

# tencentcloud_lighthouse_start_instance

Provides a resource to create a lighthouse start_instance

## Example Usage

```hcl
resource "tencentcloud_lighthouse_start_instance" "start_instance" {
  instance_id = "lhins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



