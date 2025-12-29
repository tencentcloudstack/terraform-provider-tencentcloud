Provides a resource to create a VCube renew video operation

~> **NOTE:** Resource `tencentcloud_vcube_renew_video_operation` can be directly invoked to renew the license within 30 days before its expiration. Once the renewal is successful, an additional year will be added immediately.

Example Usage

```hcl
resource "tencentcloud_vcube_renew_video_operation" "example" {
  license_id = 1513
}
```
