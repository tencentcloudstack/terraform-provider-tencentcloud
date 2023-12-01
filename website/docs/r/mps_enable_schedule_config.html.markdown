---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_enable_schedule_config"
sidebar_current: "docs-tencentcloud-resource-mps_enable_schedule_config"
description: |-
  Provides a resource to create a mps enable_schedule_config
---

# tencentcloud_mps_enable_schedule_config

Provides a resource to create a mps enable_schedule_config

## Example Usage

### Enable the mps schedule

```hcl
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-schedule-config-output1-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_mps_schedule" "example" {
  schedule_name = "tf_test_mps_schedule_config"

  trigger {
    type = "CosFileUpload"
    cos_file_upload_trigger {
      bucket  = data.tencentcloud_cos_bucket_object.object.bucket
      region  = "%s"
      dir     = "/upload/"
      formats = ["flv", "mov"]
    }
  }

  activities {
    activity_type   = "input"
    reardrive_index = [1, 2]
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [3]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [6, 7]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [4, 5]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [10]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [10]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [10]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [8]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [9]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [10]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type = "output"
  }

  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }

  output_dir = "output/"
}

resource "tencentcloud_mps_enable_schedule_config" "config" {
  schedule_id = tencentcloud_mps_schedule.example.id
  enabled     = true
}
```

### Disable the mps schedule

```hcl
resource "tencentcloud_mps_enable_schedule_config" "config" {
  schedule_id = tencentcloud_mps_schedule.example.id
  enabled     = false
}
```

## Argument Reference

The following arguments are supported:

* `enabled` - (Required, Bool) true: enable; false: disable.
* `schedule_id` - (Required, Int) The scheme ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps enable_schedule_config can be imported using the id, e.g.

```
terraform import tencentcloud_mps_enable_schedule_config.enable_schedule_config enable_schedule_config_id
```

