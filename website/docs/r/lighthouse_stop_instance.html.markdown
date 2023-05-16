---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_stop_instance"
sidebar_current: "docs-tencentcloud-resource-lighthouse_stop_instance"
description: |-
  Provides a resource to create a lighthouse stop_instance
---

# tencentcloud_lighthouse_stop_instance

Provides a resource to create a lighthouse stop_instance

## Example Usage

```hcl
resource "tencentcloud_lighthouse_stop_instance" "stop_instance" {
  instance_id = "lhins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



