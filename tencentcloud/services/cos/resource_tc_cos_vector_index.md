Provides a resource to create a COS vector index.

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

resource "tencentcloud_cos_vector_index" "example" {
  vector_bucket_name = tencentcloud_cos_vector_bucket.example.id
  index_name         = "tf-example"
  data_type          = "float32"
  dimension          = 1
  distance_metric    = "euclidean"

  metadata_configuration = {
    non_filterable_metadata_keys = [
      "key1",
      "key2",
      "key3"
    ]
  }
}
```

Import

COS vector index can be imported using the composite id vector_bucket_name#index_name, e.g.

```
terraform import tencentcloud_cos_vector_index.example tf-example-1309101522#tf-example
```
