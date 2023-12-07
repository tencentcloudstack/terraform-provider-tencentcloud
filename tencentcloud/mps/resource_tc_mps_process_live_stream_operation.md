Provides a resource to create a mps process_live_stream_operation

Example Usage

Process mps live stream through CMQ

```hcl
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-process-live-stream-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_mps_process_live_stream_operation" "operation" {
  url = "http://www.abc.com/abc.m3u8"
  task_notify_config {
    cmq_model   = "Queue"
    cmq_region  = "gz"
    queue_name  = "test"
    topic_name  = "test"
    notify_type = "CMQ"
  }

  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }

  output_dir = "/output/"

  ai_content_review_task {
    definition = 10
  }

  ai_recognition_task {
    definition = 10
  }
}
```