Provides a resource to create a tsf application_public_config_release

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config_release" "application_public_config_release" {
  config_id = "dcfg-p-123456"
  namespace_id = "namespace-123456"
  release_desc = "product version"
}
```

Import

tsf application_public_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_public_config_release.application_public_config_release application_public_config_attachment_id
```