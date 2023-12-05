Provides a resource to create a cam set_policy_version_config

Example Usage

```hcl
resource "tencentcloud_cam_set_policy_version_config" "set_policy_version_config" {
  policy_id = 171162811
  version_id = 2
}
```

Import

cam set_policy_version_config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_set_policy_version_config.set_policy_version_config set_policy_version_config_id
```