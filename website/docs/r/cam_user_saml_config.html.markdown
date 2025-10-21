---
subcategory: "Cloud Access Management(CAM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cam_user_saml_config"
sidebar_current: "docs-tencentcloud-resource-cam_user_saml_config"
description: |-
  Provides a resource to create a CAM user saml config
---

# tencentcloud_cam_user_saml_config

Provides a resource to create a CAM user saml config

## Example Usage

### Create saml config by metadata string

```hcl
resource "tencentcloud_cam_user_saml_config" "example" {
  saml_metadata_document = <<-EOT
  <?xml version="1.0" encoding="utf-8"?></EntityDescriptor>
EOT
  auxiliary_domain       = "xxx.com"
}
```

### Create saml config by metadata file path

```hcl
resource "tencentcloud_cam_user_saml_config" "example" {
  saml_metadata_document = "./metadataDocument.xml"
}
```

## Argument Reference

The following arguments are supported:

* `saml_metadata_document` - (Required, String) SAML metadata document, xml format, support string content or file path.
* `auxiliary_domain` - (Optional, String) auxiliary domain, like: `xxx.com`.
* `metadata_document_file` - (Optional, String) The path used to save the saml Metadata file.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status: `0`: not set, `11`: enabled, `2`: disabled.


## Import

CAM user saml config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_saml_config.example 79f23f0f-ad00-414f-bc5e-94859ffdfb9e
```

