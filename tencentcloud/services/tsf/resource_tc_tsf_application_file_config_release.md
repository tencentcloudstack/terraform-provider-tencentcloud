Provides a resource to create a tsf application_file_config_release

Example Usage

```hcl
resource "tencentcloud_tsf_application_file_config_release" "application_file_config_release" {
  config_id = "dcfg-f-123456"
  group_id = "group-123456"
  release_desc = "product release"
}
```

Import

tsf applicationfile_config_release can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_application_file_config_release.application_file_config_release application_file_config_release_id
```