package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsSchedulesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsSchedulesDataSource_none,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.activities.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.update_time"),
				),
			},
			{
				Config: fmt.Sprintf(testAccMpsSchedulesDataSource_specific_one, defaultMpsScheduleId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_ids.#"),
					resource.TestCheckTypeSetElemAttr("data.tencentcloud_mps_schedules.schedules", "schedule_ids.*", fmt.Sprint(defaultMpsScheduleId)),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "trigger_type", "CosFileUpload"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "status", "Disabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_id", fmt.Sprint(defaultMpsScheduleId)),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_name", defaultMpsScheduleName),
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

const testAccMpsSchedulesDataSource_specific_one = `

data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = [%d]
  trigger_type = "CosFileUpload"
  status = "Disabled"
  }

`
