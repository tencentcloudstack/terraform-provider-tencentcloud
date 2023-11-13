package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMongodbInstanceSlowLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstanceSlowLogDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_slow_log.instance_slow_log")),
			},
		},
	})
}

const testAccMongodbInstanceSlowLogDataSource = `

data "tencentcloud_mongodb_instance_slow_log" "instance_slow_log" {
  instance_id = "cmgo-9d0p6umb"
  start_time = "2019-06-01 10:00:00"
  end_time = "2019-06-02 12:00:00"
  slow_m_s = 100
  format = "json"
  }

`
