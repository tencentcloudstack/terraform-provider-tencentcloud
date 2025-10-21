---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_disk_attachment"
sidebar_current: "docs-tencentcloud-resource-lighthouse_disk_attachment"
description: |-
  Provides a resource to create a lighthouse disk_attachment
---

# tencentcloud_lighthouse_disk_attachment

Provides a resource to create a lighthouse disk_attachment

## Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_attachment" "disk_attachment" {
  disk_id     = "lhdisk-xxxxxx"
  instance_id = "lhins-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, String, ForceNew) Disk id.
* `instance_id` - (Required, String, ForceNew) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

lighthouse disk_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_disk_attachment.disk_attachment disk_attachment_id
```

