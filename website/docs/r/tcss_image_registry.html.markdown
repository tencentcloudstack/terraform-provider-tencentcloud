---
subcategory: "Tencent Container Security Service(TCSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcss_image_registry"
sidebar_current: "docs-tencentcloud-resource-tcss_image_registry"
description: |-
  Provides a resource to create a tcss image registry
---

# tencentcloud_tcss_image_registry

Provides a resource to create a tcss image registry

## Example Usage

```hcl
resource "tencentcloud_tcss_image_registry" "example" {
  name             = "terraform"
  username         = "root"
  password         = "Password@demo"
  url              = "https://example.com"
  registry_type    = "harbor"
  net_type         = "public"
  registry_version = "V1"
  registry_region  = "default"
  need_scan        = true
  conn_detect_config {
    quuid = "backend"
    uuid  = "backend"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Repository name.
* `net_type` - (Required, String) Network type, which can be `public` (public network).
* `password` - (Required, String) Password.
* `registry_type` - (Required, String) Repository type, which can be `harbor`. Valid values: harbor, quay, jfrog, aws, azure, other-tcr.
* `url` - (Required, String) Repository URL.
* `username` - (Required, String) Username.
* `conn_detect_config` - (Optional, Set) Connectivity detection configuration.
* `insecure` - (Optional, Int) Valid values: `0` (secure mode with certificate verification, which is the default value); `1` (unsecure mode that skips certificate verification).
* `need_scan` - (Optional, Bool) Whether to scan the latest image.
* `registry_region` - (Optional, String) Region. Default value: `default`.
* `registry_version` - (Optional, String) Repository version.
* `speed_limit` - (Optional, Int) Speed limit.

The `conn_detect_config` object supports the following:

* `quuid` - (Optional, String) Host Quuid.
* `uuid` - (Optional, String) Host uuid.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `sync_status` - Sync status.


