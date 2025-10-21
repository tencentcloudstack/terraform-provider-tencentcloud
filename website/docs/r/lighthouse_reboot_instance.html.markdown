---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_reboot_instance"
sidebar_current: "docs-tencentcloud-resource-lighthouse_reboot_instance"
description: |-
  Provides a resource to create a lighthouse reboot_instance
---

# tencentcloud_lighthouse_reboot_instance

Provides a resource to create a lighthouse reboot_instance

## Example Usage

```hcl
resource "tencentcloud_lighthouse_reboot_instance" "reboot_instance" {
  instance_id = "lhins-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



