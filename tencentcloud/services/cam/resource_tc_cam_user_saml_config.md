Provides a resource to create a CAM user saml config

Example Usage

Create saml config by metadata string

```hcl
resource "tencentcloud_cam_user_saml_config" "example" {
  saml_metadata_document = <<-EOT
  <?xml version="1.0" encoding="utf-8"?></EntityDescriptor>
EOT
  auxiliary_domain       = "xxx.com"
}
```

Create saml config by metadata file path

```hcl
resource "tencentcloud_cam_user_saml_config" "example" {
  saml_metadata_document = "./metadataDocument.xml"
}
```

Import

CAM user saml config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_saml_config.example 79f23f0f-ad00-414f-bc5e-94859ffdfb9e
```