---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_images"
sidebar_current: "docs-tencentcloud-datasource-tcr_images"
description: |-
  Use this data source to query detailed information of tcr images
---

# tencentcloud_tcr_images

Use this data source to query detailed information of tcr images

## Example Usage

```hcl
data "tencentcloud_tcr_images" "images" {
  registry_id     = "tcr-xxx"
  namespace_name  = "ns"
  repository_name = "repo"
  image_version   = "v1"
  digest          = "sha256:xxxxx"
  exact_match     = false
}
```

## Argument Reference

The following arguments are supported:

* `namespace_name` - (Required, String) namespace name.
* `registry_id` - (Required, String) instance id.
* `repository_name` - (Required, String) repository name.
* `digest` - (Optional, String) specify image digest for lookup.
* `exact_match` - (Optional, Bool) specifies whether it is an exact match, true is an exact match, and not filled is a fuzzy match.
* `image_version` - (Optional, String) image version name, default is fuzzy match.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `image_info_list` - container image information list.
  * `digest` - hash value.
  * `image_version` - tag name.
  * `kind` - product type,note: this field may return null, indicating that no valid value can be obtained.
  * `kms_signature` - kms signature information,note: this field may return null, indicating that no valid value can be obtained.
  * `size` - image size (unit: byte).
  * `update_time` - update time.


