resource "tencentcloud_cos_bucket" "bucket" {
  bucket = "${var.bucket-name}"
  acl    = "${var.acl}"

  cors_rules {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "POST"]
    allowed_origins = ["https://www.test.com"]
    expose_headers  = ["x-cos-test"]
    max_age_seconds = 300
  }

  lifecycle_rules {
    filter_prefix = "test/"

    expiration {
      days = 365
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "ARCHIVE"
    }
  }

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}

resource "tencentcloud_cos_bucket_object" "object" {
  bucket       = "${tencentcloud_cos_bucket.bucket.bucket}"
  key          = "${var.object-name}"
  content      = "${var.object-content}"
  content_type = "binary/octet-stream"
}

data "tencentcloud_cos_bucket_object" "data_object" {
  bucket = "${tencentcloud_cos_bucket.bucket.id}"
  key    = "${tencentcloud_cos_bucket_object.object.key}"
}

data "tencentcloud_cos_buckets" "data_bucket" {
  bucket_prefix = "${tencentcloud_cos_bucket.bucket.id}"
}
