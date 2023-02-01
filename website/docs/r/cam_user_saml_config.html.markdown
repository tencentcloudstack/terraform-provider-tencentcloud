---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user_saml_config"
sidebar_current: "docs-tencentcloud-resource-cam_user_saml_config"
description: |-
  Provides a resource to create a cam user_saml_config
---

# tencentcloud_cam_user_saml_config

Provides a resource to create a cam user_saml_config

## Example Usage

```hcl
resource "tencentcloud_cam_user_saml_config" "user_saml_config" {
  saml_metadata_document = "./metadataDocument.xml"
  # saml_metadata_document  = <<-EOT
  # <?xml version="1.0" encoding="utf-8"?></EntityDescriptor>
  # EOT
}
```

## Argument Reference

The following arguments are supported:

* `saml_metadata_document` - (Required, String) SAML metadata document, xml format, support string content or file path.
* `metadata_document_file` - (Optional, String) The path used to save the samlMetadat file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status: `0`: not set, `11`: enabled, `2`: disabled.


## Import

cam user_saml_config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_saml_config.user_saml_config user_id
```

