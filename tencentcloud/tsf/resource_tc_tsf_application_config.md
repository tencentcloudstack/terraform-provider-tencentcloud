Provides a resource to create a tsf application_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_config" "application_config" {
  config_name = "test-2"
  config_version = "1.0"
  config_value = "name: \"name\""
  application_id = "application-ym9mxmza"
  config_version_desc = "test2"
  # config_type = ""
  encode_with_base64 = false
  # program_id_list =
}
```