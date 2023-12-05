Provides a resource to create a cam service_linked_role

Example Usage

```hcl
resource "tencentcloud_cam_service_linked_role" "service_linked_role" {
  qcs_service_name = ["cvm.qcloud.com","ekslog.tke.cloud.tencent.com"]
  custom_suffix = "tf"
  description = "desc cam"
  tags = {
    "createdBy" = "terraform"
  }
}

```