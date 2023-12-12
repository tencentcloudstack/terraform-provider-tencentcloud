Provides a resource to create a cam user_saml_config

Example Usage

```hcl
resource "tencentcloud_cam_user_saml_config" "user_saml_config" {
  saml_metadata_document = "./metadataDocument.xml"
  # saml_metadata_document  = <<-EOT
  # <?xml version="1.0" encoding="utf-8"?></EntityDescriptor>
  # EOT
}
```

Import

cam user_saml_config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_user_saml_config.user_saml_config user_id
```