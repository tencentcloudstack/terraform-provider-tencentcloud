Provides a resource to create a tsf api_group

Example Usage

```hcl
resource "tencentcloud_tsf_api_group" "api_group" {
  group_name = "terraform_test_group"
  group_context = "/terraform-test"
  auth_type = "none"
  description = "terraform-test"
  group_type = "ms"
  gateway_instance_id = "gw-ins-i6mjpgm8"
  # namespace_name_key = "path"
  # service_name_key = "path"
  namespace_name_key_position = "path"
  service_name_key_position = "path"
}
```

Import

tsf api_group can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_group.api_group api_group_id
```