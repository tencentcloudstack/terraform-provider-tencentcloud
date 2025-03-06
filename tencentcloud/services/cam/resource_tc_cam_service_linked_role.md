Provides a resource to create a CAM service linked role

Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "example" {
  qcs_service_name = ["cvm.qcloud.com", "ekslog.tke.cloud.tencent.com"]
  custom_suffix    = "tf-example"
  description      = "description."
  tags = {
    createdBy = "Terraform"
  }
}
```

Import

CAM service linked role can be imported using the id, e.g.

```
$ terraform import tencentcloud_cam_service_linked_role.example 4611686018441982195
```