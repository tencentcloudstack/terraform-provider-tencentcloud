---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_image_manifests"
sidebar_current: "docs-tencentcloud-datasource-tcr_image_manifests"
description: |-
  Use this data source to query detailed information of tcr image_manifests
---

# tencentcloud_tcr_image_manifests

Use this data source to query detailed information of tcr image_manifests

## Example Usage

```hcl
data "tencentcloud_tcr_image_manifests" "image_manifests" {
  registry_id     = "%s"
  namespace_name  = "%s"
  repository_name = "%s"
  image_version   = "v1"
}
```

## Argument Reference

The following arguments are supported:

* `image_version` - (Required, String) mirror version.
* `namespace_name` - (Required, String) namespace name.
* `registry_id` - (Required, String) instance ID.
* `repository_name` - (Required, String) mirror warehouse name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `config` - configuration information of the image.
* `manifest` - Manifest information of the image.


