Provides a resource to manage a COS vector bucket.

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_vector_bucket" "example" {
  vector_bucket_name = "tf-example-${local.app_id}"

  encryption_configuration = {
    sse_type = "AES256"
  }
}
```

Import

COS vector bucket can be imported using the id, e.g.

```
terraform import tencentcloud_cos_vector_bucket.example tf-example-1309101522
```
