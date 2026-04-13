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

## Argument Reference

The following arguments are supported:

* `input_info` - (Required) The information of the file to process.
* `output_storage` - (Optional) The output storage of media processing.
* `output_dir` - (Optional) The output directory of media processing.
* `ai_content_review_task` - (Optional) The AI content review task configuration.
* `ai_recognition_task` - (Optional) The AI recognition task configuration.
* `task_notify_config` - (Optional) The task notification configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_id` - Task ID returned by API, used to track media processing task status.
* `id` - Resource ID, same as `task_id`.