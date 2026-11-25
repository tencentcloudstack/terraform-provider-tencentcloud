---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_index"
sidebar_current: "docs-tencentcloud-resource-cos_vector_index"
description: |-
  Provides a resource to create a COS vector index.
---

# tencentcloud_cos_vector_index

Provides a resource to create a COS vector index.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `data_type` - (Required, String) Data type of the vectors. Default value: `float32`. It cannot be changed after creation.
* `dimension` - (Required, Int) Dimension of the vectors. Value range: 1-4096. It cannot be changed after creation.
* `distance_metric` - (Required, String) Distance metric for similarity search. Valid values: `euclidean`, `cosine`. It cannot be changed after creation.
* `index_name` - (Required, String) Name of the vector index, composed of lowercase letters, numbers, `-` and `.`, beginning and ending with a lowercase letter or number, 3-63 characters in length. It cannot be changed after creation.
* `vector_bucket_name` - (Required, String) Name of the vector bucket the index belongs to, in the format of `<BucketName-APPID>`, e.g. `examplebucket-1250000000`. It cannot be changed after creation.
* `metadata_configuration` - (Optional, Object) Metadata configuration of the index. It cannot be changed after creation.

The `metadata_configuration` object supports the following:

* `non_filterable_metadata_keys` - (Required, List) Metadata keys that cannot be used for filtering. At least 1 and at most 10 keys, each 1-63 characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource, in the format of `vector_bucket_name#index_name`.


## Import

COS vector index can be imported using the composite id vector_bucket_name#index_name, e.g.

```
terraform import tencentcloud_cos_vector_index.example tf-example-1309101522#tf-example
```

