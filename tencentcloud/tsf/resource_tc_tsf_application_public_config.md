Provides a resource to create a tsf application_public_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_public_config" "application_public_config" {
  config_name = "my_config"
  config_version = "1.0"
  config_value = "test: 1"
  config_version_desc = "product version"
  config_type = "P"
  encode_with_base64 = true
  # program_id_list =
}
```