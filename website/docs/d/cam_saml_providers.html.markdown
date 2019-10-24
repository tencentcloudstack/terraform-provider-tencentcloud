---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_saml_providers"
sidebar_current: "docs-tencentcloud-datasource-cam_saml_providers"
description: |-
  Use this data source to query detailed information of CAM SAML providers
---

# tencentcloud_cam_saml_providers

Use this data source to query detailed information of CAM SAML providers

## Example Usage

```hcl
data "tencentcloud_cam_saml_providers" "foo" {
  name = "cam-test-provider"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the CAM SAML provider .
* `name` - (Optional) Name of the CAM SAML provider to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `provider_list` - A list of CAM SAML providers. Each element contains the following attributes:
  * `create_time` - Create time of the CAM SAML provider.
  * `description` - Description of CAM SAML provider.
  * `modify_time` - The last modify time of the CAM SAML provider.
  * `name` - Name of CAM SAML provider.


