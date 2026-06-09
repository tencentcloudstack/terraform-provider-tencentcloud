---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_check_cname_status_operation"
sidebar_current: "docs-tencentcloud-resource-teo_check_cname_status_operation"
description: |-
  Provides a resource to check CNAME status for TEO domains.
---

# tencentcloud_teo_check_cname_status_operation

Provides a resource to check CNAME status for TEO domains.

## Example Usage

```hcl
resource "tencentcloud_teo_check_cname_status_operation" "example" {
  zone_id = "zone-12345678"
  record_names = [
    "example.com",
    "www.example.com",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `record_names` - (Required, List: [`String`], ForceNew) List of record names to check CNAME status.
* `zone_id` - (Required, String, ForceNew) Zone ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cname_status` - CNAME status information for each record name.
  * `cname` - CNAME address. May be null.
  * `record_name` - Record name.
  * `status` - CNAME status. Valid values: active, moved.


