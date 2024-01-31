package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsSchedulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSchedulesDataSource_none,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.activities.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.update_time"),
				),
			},
			{
				Config: testAccMpsSchedulesDataSource_specific_one,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_ids.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "trigger_type", "CosFileUpload"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "status", "Disabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_name", "tf_test_mps_schedule_001"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.status", "Disabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.trigger.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.activities.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.output_storage.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.output_dir"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.task_notify_config.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.update_time"),
				),
			},
		},
	})
}

const testAccMpsSchedulesDataSource_none = `

data "tencentcloud_mps_schedules" "schedules" {
  }

`

const testAccMpsSchedulesDataSource_specific_one = tcacctest.UserInfoData + `
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-schedule-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_cos_bucket_object" "object" {
  bucket       = tencentcloud_cos_bucket.output.bucket
  key          = "/mps-test/test.mov"
  content      = "aaaaaaaaaaaaaaaa"
  content_type = "binary/octet-stream"
}

resource "tencentcloud_mps_schedule" "schedule" {
  schedule_name = "tf_test_mps_schedule_001"

  trigger {
    type = "CosFileUpload"
    cos_file_upload_trigger {
      bucket  = tencentcloud_cos_bucket_object.object.bucket
      region  = "ap-guangzhou"
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
      region = "ap-guangzhou"
    }
  }

  output_dir = "output/"
  resource_id = "vts-2600014161-0"
}

data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = [tencentcloud_mps_schedule.schedule.id]
  trigger_type = "CosFileUpload"
  status = "Disabled"
}

`
