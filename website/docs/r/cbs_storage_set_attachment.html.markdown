---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storage_set_attachment"
sidebar_current: "docs-tencentcloud-resource-cbs_storage_set_attachment"
description: |-
  Provides a CBS storage set attachment resource.
---

# tencentcloud_cbs_storage_set_attachment

Provides a CBS storage set attachment resource.

## Example Usage

```hcl
resource "tencentcloud_cbs_storage_set_attachment" "attachment" {
  storage_id  = "disk-kdt0sq6m"
  instance_id = "ins-jqlegd42"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of the CVM instance.
* `storage_id` - (Required, String, ForceNew) ID of the mounted CBS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



