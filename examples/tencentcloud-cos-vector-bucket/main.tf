resource "tencentcloud_cos_vector_bucket" "example" {
  vector_bucket_name = var.vector_bucket_name
  force_destroy      = false

  encryption_configuration {
    sse_type = "AES256"
  }
}

resource "tencentcloud_cos_vector_index" "example" {
  vector_bucket_name = tencentcloud_cos_vector_bucket.example.vector_bucket_name
  index_name         = var.index_name
  data_type          = "float32"
  dimension          = var.dimension
  distance_metric    = var.distance_metric

  metadata_configuration {
    non_filterable_metadata_keys = ["source_text"]
  }
}

resource "tencentcloud_cos_vector_bucket_policy" "example" {
  count             = var.policy_principal_qcs == null ? 0 : 1
  vector_bucket_qcs = tencentcloud_cos_vector_bucket.example.vector_bucket_qcs
  policy = jsonencode({
    version = "2.0"
    statement = [{
      effect    = "allow"
      principal = { qcs = [var.policy_principal_qcs] }
      action    = ["name/cos:QueryVectors"]
      resource  = [tencentcloud_cos_vector_index.example.index_qcs]
    }]
  })
}
