Provides a resource to create a mps enable_workflow_config

Example Usage

Enable the mps workflow

```hcl
resource "tencentcloud_mps_workflow" "example" {
  output_dir    = "/"
  task_priority = 0
  workflow_name = "tf-workflow-enable-config"

  media_process_task {
    adaptive_dynamic_streaming_task_set {
      definition             = 12
      output_object_path     = "/out"
      segment_object_name    = "/out"
      sub_stream_object_name = "/out/out/"

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    snapshot_by_time_offset_task_set {
      definition          = 10
      ext_time_offset_set = [
        "1s",
      ]
      output_object_path  = "/snapshot/"
      time_offset_set     = []

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }

    animated_graphic_task_set {
      definition         = 20000
      end_time_offset    = 0
      output_object_path = "/test/"
      start_time_offset  = 0

      output_storage {
        type = "COS"

        cos_output_storage {
          bucket = "cos-lock-1308919341"
          region = "ap-guangzhou"
        }
      }
    }
  }

  ai_analysis_task {
    definition = 20
  }

  ai_content_review_task {
    definition = 20
  }

  ai_recognition_task {
    definition = 20
  }

  output_storage {
    type = "COS"

    cos_output_storage {
      bucket = "cos-lock-1308919341"
      region = "ap-guangzhou"
    }
  }

  trigger {
    type = "CosFileUpload"

    cos_file_upload_trigger {
      bucket = "cos-lock-1308919341"
      dir    = "/"
      region = "ap-guangzhou"
    }
  }
}

resource "tencentcloud_mps_enable_workflow_config" "config" {
  workflow_id = tencentcloud_mps_workflow.example.id
  enabled = true
}

```

Disable the mps workflow

```hcl
resource "tencentcloud_mps_enable_workflow_config" "config" {
  workflow_id = tencentcloud_mps_workflow.example.id
  enabled = false
}

```

Import

mps enable_workflow_config can be imported using the id, e.g.

```
terraform import tencentcloud_mps_enable_workflow_config.enable_workflow_config enable_workflow_config_id
```