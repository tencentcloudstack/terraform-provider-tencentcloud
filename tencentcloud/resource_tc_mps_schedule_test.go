package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsScheduleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsSchedule, defaultRegion, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "schedule_name", "tf_test_mps_schedule"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.type", "CosFileUpload"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.dir", "/upload/"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "mp4"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "mov"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "input"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "2"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "3"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.0.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "6"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "7"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.0.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "output"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.region", defaultRegion),

					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage", "output/"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_schedule.schedule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsSchedule = userInfoData + `
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-schedule-output-${local.app_id}"
  force_clean = true
}

resource "tencentcloud_mps_schedule" "schedule" {
  schedule_name = "tf_test_mps_schedule"

  trigger {
    type = "CosFileUpload"
    cos_file_upload_trigger {
      bucket  = data.tencentcloud_cos_bucket_object.object.bucket
      region  = "%s"
      dir     = "/upload/"
      formats = ["mp4", "mov"]
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


`
