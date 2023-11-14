package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbSlowLogDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbSlowLogDataDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_slow_log_data.slow_log_data")),
			},
		},
	})
}

const testAccCdbSlowLogDataDataSource = `

data "tencentcloud_cdb_slow_log_data" "slow_log_data" {
  instance_id = ""
  start_time = 
  end_time = 
  user_hosts = 
  user_names = 
  data_bases = 
  sort_by = ""
  order_by = ""
  inst_type = ""
  }

`
