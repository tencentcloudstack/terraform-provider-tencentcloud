package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOceanusLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusLogDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_log.log")),
			},
		},
	})
}

const testAccOceanusLogDataSource = `

data "tencentcloud_oceanus_log" "log" {
  job_id = "cql-6v1jkxrn"
  start_time = 1611754219108
  end_time = 1611754219108
  running_order_id = 1
  keyword = "xx"
    order_type = "asc"
          }

`
