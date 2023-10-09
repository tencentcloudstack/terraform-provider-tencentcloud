package tencentcloud

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("tencentcloud_mps_schedule", &resource.Sweeper{
		Name: "tencentcloud_mps_schedule",
		F:    testSweepMpsSchedule,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mps_schedule
func testSweepMpsSchedule(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	mpsService := MpsService{client: cli.(*TencentCloudClient).apiV3Conn}

	schedules, err := mpsService.DescribeMpsScheduleById(ctx, nil)
	if err != nil {
		return err
	}
	if len(schedules) == 0 {
		return fmt.Errorf("mps schedules not exists.")
	}

	for _, v := range schedules {
		delId := *v.ScheduleId
		delName := *v.ScheduleName

		if strings.HasPrefix(delName, "tf_test_") {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := mpsService.DeleteMpsScheduleById(ctx, helper.Int64ToStr(delId))
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] delete mps schedule %d failed! reason:[%s]", delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudMpsScheduleResource_basic(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsSchedule, randomNum, defaultRegion, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "schedule_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.type", "CosFileUpload"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.dir", "/upload/"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "flv"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "mov"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "input"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "2"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.1.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.1.reardrive_index.*", "3"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.1.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.1.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.2.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.2.reardrive_index.*", "6"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.2.reardrive_index.*", "7"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.2.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.2.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.3.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.3.reardrive_index.*", "4"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.3.reardrive_index.*", "5"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.3.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.3.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.4.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.4.reardrive_index.*", "10"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.4.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.4.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.10.activity_type", "output"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.region", defaultRegion),

					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_dir", "output/"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsSchedule_update, randomNum, defaultRegion, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "schedule_name", fmt.Sprintf("tf_test_mps_schedule_%d_changed", randomNum)),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.type", "CosFileUpload"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.dir", "/upload_changed/"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "mp4"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "trigger.0.cos_file_upload_trigger.0.formats.*", "mov"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.0.activity_type", "input"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.0.reardrive_index.*", "2"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.1.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.1.reardrive_index.*", "3"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.1.reardrive_index.*", "4"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.1.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.1.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.2.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.2.reardrive_index.*", "7"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.2.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.2.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.3.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.3.reardrive_index.*", "5"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.3.reardrive_index.*", "6"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.3.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.3.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.4.activity_type", "action-trans"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_schedule.schedule", "activities.4.reardrive_index.*", "10"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.4.activity_para.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.4.activity_para.0.transcode_task.0.definition", "10"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "activities.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "activities.10.activity_type", "output"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_storage.0.cos_output_storage.0.region", defaultRegion),

					resource.TestCheckResourceAttr("tencentcloud_mps_schedule.schedule", "output_dir", "output_chagned/"),
				),
			},
			{
				ResourceName: "tencentcloud_mps_schedule.schedule",
				ImportState:  true,
				// ImportStateVerify: true,
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
  acl         = "public-read"
}

resource "tencentcloud_mps_schedule" "schedule" {
  schedule_name = "tf_test_mps_schedule_%d"

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

const testAccMpsSchedule_update = userInfoData + `
data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-schedule-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_mps_schedule" "schedule" {
  schedule_name = "tf_test_mps_schedule_%d_changed"

  trigger {
    type = "CosFileUpload"
    cos_file_upload_trigger {
      bucket  = data.tencentcloud_cos_bucket_object.object.bucket
      region  = "%s"
      dir     = "/upload_changed/"
      formats = ["mp4", "mov"]
    }
  }

  activities {
    activity_type   = "input"
    reardrive_index = [1, 2]
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [3, 4]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [7]
    activity_para {
      transcode_task {
        definition = 10
      }
    }
  }

  activities {
    activity_type   = "action-trans"
    reardrive_index = [5, 6]
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

  output_dir = "output_chagned/"
}


`
