Provides a resource to create a cos bucket inventory

Example Usage

```hcl
# get user info
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

# create cos
resource "tencentcloud_cos_bucket" "example" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}

# create cos bucket inventory
resource "tencentcloud_cos_bucket_inventory" "example" {
  name                     = "tf-example"
  bucket                   = tencentcloud_cos_bucket.example.id
  is_enabled               = "true"
  included_object_versions = "Current"

  optional_fields {
    fields = ["Size", "ETag"]
  }

  filter {
    period {
      start_time = "1687276800"
    }
  }

  schedule {
    frequency = "Daily"
  }

  destination {
    bucket = "qcs::cos:ap-guangzhou::private-bucket-1309118522"
    format = "CSV"
    prefix = "frontends"
  }
}
```

Import

cos bucket inventory can be imported using the id, e.g.

```
terraform import tencentcloud_cos_bucket_inventory.example private-bucket-1309118522#tf-example
```
