package mps_test

import (
	"fmt"
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
				Config: fmt.Sprintf(testAccMpsSchedulesDataSource_specific_one, tcacctest.DefaultMpsScheduleId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_ids.#"),
					resource.TestCheckTypeSetElemAttr("data.tencentcloud_mps_schedules.schedules", "schedule_ids.*", fmt.Sprint(tcacctest.DefaultMpsScheduleId)),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "trigger_type", "CosFileUpload"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "status", "Disabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_id", fmt.Sprint(tcacctest.DefaultMpsScheduleId)),
					resource.TestCheckResourceAttr("data.tencentcloud_mps_schedules.schedules", "schedule_info_set.0.schedule_name", tcacctest.DefaultMpsScheduleName),
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
