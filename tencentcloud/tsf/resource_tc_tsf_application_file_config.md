Provides a resource to create a tsf application_file_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config" "application_file_config" {
  config_name = "terraform-test"
  config_version = "1.0"
  config_file_name = "application.yaml"
  config_file_value = "test: 1"
  application_id = "application-a24x29xv"
  config_file_path = "/etc/nginx"
  config_version_desc = "1.0"
  config_file_code = "UTF-8"
  config_post_cmd = "source .bashrc"
  encode_with_base64 = true
}
```