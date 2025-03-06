Provides a resource to create a CAM set policy version config

Example Usage

```hcl
resource "tencentcloud_cam_set_policy_version_config" "example" {
  policy_id  = 234290251
  version_id = 3
}
```

Import

CAM set policy version config can be imported using the id, e.g.

```
terraform import tencentcloud_cam_set_policy_version_config.example 234290251#3
```