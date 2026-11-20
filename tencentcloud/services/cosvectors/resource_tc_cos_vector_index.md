Provides a resource to manage a COS Vectors index.

Use a reference to the vector bucket name to establish dependency ordering. Vector writes and queries remain application SDK/API operations; Terraform does not store vector data. Data type is `float32`, dimension is an integer from 1 to 4096, and distance metric is `cosine` or `euclidean`. Index names must be 3-63 lowercase letters, digits, dots or hyphens, start and end with a letter or digit, and must not contain `..`.

-> **Note:** Changing the bucket, index name, data type, dimension, distance metric or explicitly configured metadata replaces the index. Deleting or replacing an index permanently deletes all its vectors, even when the bucket has `force_destroy = false`.

Replace the sample APPID with your actual APPID and configure the provider region and environment credentials. Keep secrets outside HCL. `metadata_configuration` accepts at most one block. Its keys are a set of 1-10 distinct strings matching `^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`. Order does not cause a diff. Omit the block to preserve metadata discovered during import; omission is not a request to clear it.

The resource returns actual `index_qcs`, not an AWS ARN. `creation_time` preserves the raw SDK integer without inferring units. Region comes from the provider; use provider aliases for multiple regions. Index-level encryption, tags, KMS, resource-level region overrides and custom timeouts are not exposed.

API

This resource uses `CreateIndex`, `GetIndex` and `DeleteIndex`.

Example Usage

```hcl
resource "tencentcloud_cos_vector_bucket" "example" {
  vector_bucket_name = "tf-vector-example-1250000000"
}

resource "tencentcloud_cos_vector_index" "example" {
  vector_bucket_name = tencentcloud_cos_vector_bucket.example.vector_bucket_name
  index_name         = "example-index"
  data_type          = "float32"
  dimension          = 3
  distance_metric    = "cosine"

  metadata_configuration {
    non_filterable_metadata_keys = ["source_text"]
  }
}
```

Import

Import an existing index using the full bucket name and index name:

```bash
terraform import tencentcloud_cos_vector_index.example tf-vector-example-1250000000/example-index
```

An index QCS is also accepted when its region matches the provider and account matches the bucket APPID:

```hcl
import {
  to = tencentcloud_cos_vector_index.example
  id = "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-example-1250000000/index/example-index"
}
```

Define matching resource configuration before import. Inspect the subsequent plan and do not apply a replacement merely to make an existing index match the example. AWS ARNs and AWS provider state cannot be imported as COS resources.
