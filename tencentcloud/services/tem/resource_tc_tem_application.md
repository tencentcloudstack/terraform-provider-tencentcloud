Provides a resource to create a tem application

Example Usage

```hcl
resource "tencentcloud_tem_application" "application" {
  application_name = "demo"
  description = "demo for test"
  coding_language = "JAVA"
  use_default_image_service = 0
  repo_type = 2
  repo_name = "qcloud/nginx"
  repo_server = "ccr.ccs.tencentyun.com"
  tags = {
    "created" = "terraform"
  }
}
```