Provides a resource to create a mps process_media_operation

Example Usage

Process mps media through CMQ

```hcl
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-edit-media-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_process_media_operation" "operation" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = data.tencentcloud_cos_bucket_object.object.bucket
      region = "%s"
      object = data.tencentcloud_cos_bucket_object.object.key
    }
  }
  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }
  output_dir = "output/"

  ai_content_review_task {
    definition = 10
  }

  ai_recognition_task {
    definition = 10
  }

  task_notify_config {
    cmq_model   = "Queue"
    cmq_region  = "gz"
    queue_name  = "test"
    topic_name  = "test"
    notify_type = "CMQ"
  }
}
```