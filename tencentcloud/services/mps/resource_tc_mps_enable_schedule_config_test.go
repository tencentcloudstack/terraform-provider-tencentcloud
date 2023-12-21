package mps_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsEnableScheduleConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsEnableScheduleConfig_enable, tcacctest.DefaultRegion, tcacctest.DefaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_schedule_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_schedule_config.config", "schedule_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_enable_schedule_config.config", "enabled", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsEnableScheduleConfig_disable, tcacctest.DefaultRegion, tcacctest.DefaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_schedule_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_schedule_config.config", "schedule_id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_enable_schedule_config.config", "enabled", "false"),
				),
			},
			{
				ResourceName: "tencentcloud_mps_enable_schedule_config.config",
				ImportState:  true,
			},
		},
	})
}

const testAccMpsSchedule_basic = tcacctest.UserInfoData + `
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
`

const testAccMpsEnableScheduleConfig_enable = testAccMpsSchedule_basic + `

resource "tencentcloud_mps_enable_schedule_config" "config" {
  schedule_id = tencentcloud_mps_schedule.example.id
  enabled = true
}

`

const testAccMpsEnableScheduleConfig_disable = testAccMpsSchedule_basic + `

resource "tencentcloud_mps_enable_schedule_config" "config" {
  schedule_id = tencentcloud_mps_schedule.example.id
  enabled = false
}

`
