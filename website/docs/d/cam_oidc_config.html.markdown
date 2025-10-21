---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_oidc_config"
sidebar_current: "docs-tencentcloud-datasource-cam_oidc_config"
description: |-
  Use this data source to query detailed information of cam oidc_config
---

# tencentcloud_cam_oidc_config

Use this data source to query detailed information of cam oidc_config

## Example Usage

```hcl
data "tencentcloud_cam_oidc_config" "oidc_config" {
  name = "cls-kzilgv5m"
}

output "identity_key" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_key
}

output "identity_url" {
  value = data.tencentcloud_cam_oidc_config.oidc_config.identity_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `client_id` - Client ID.
* `description` - Description.
* `identity_key` - Public key for signature.
* `identity_url` - IdP URL.
* `provider_type` - IdP type. 11: Role IdP.
* `status` - Status. 0: Not set; 2: Disabled; 11: Enabled.


