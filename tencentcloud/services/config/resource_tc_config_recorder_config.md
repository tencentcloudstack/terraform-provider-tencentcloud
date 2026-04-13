Provides a resource to manage Config recorder configuration (global singleton).

Example Usage

Enable monitoring and specify resource types

```hcl
resource "tencentcloud_config_recorder_config" "example" {
  status = true
  resource_types = [
    "QCS::CAM::Group",
    "QCS::CAM::Role",
    "QCS::CAM::Policy",
    "QCS::CAM::User",
    "QCS::CVM::Instance",
    "QCS::COS::Bucket",
  ]
}
```

Disable monitoring

```hcl
resource "tencentcloud_config_recorder_config" "example" {
  status = false
}
```

Import

Config recorder config can be imported using its token ID, e.g.

```
terraform import tencentcloud_config_recorder_config.example <id>
```
