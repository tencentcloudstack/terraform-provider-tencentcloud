---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application_public_config_release"
sidebar_current: "docs-tencentcloud-resource-tsf_application_public_config_release"
description: |-
  Provides a resource to create a tsf application_public_config_release
---

# tencentcloud_tsf_application_public_config_release

Provides a resource to create a tsf application_public_config_release

## Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config_release" "application_public_config_release" {
  config_id    = "dcfg-p-123456"
  namespace_id = "namespace-123456"
  release_desc = "product version"
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required, String, ForceNew) ConfigId.
* `namespace_id` - (Required, String, ForceNew) namespace-id.
* `release_desc` - (Optional, String, ForceNew) Release description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tsf application_public_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_public_config_release.application_public_config_release application_public_config_attachment_id
```

