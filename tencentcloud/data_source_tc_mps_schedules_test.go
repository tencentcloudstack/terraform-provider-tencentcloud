package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Config: testAccMpsSchedulesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_schedules.schedules")),
			},
		},
	})
}

const testAccMpsSchedulesDataSource = `

data "tencentcloud_mps_schedules" "schedules" {
  schedule_ids = 
  trigger_type = ""
  status = ""
  }

`
