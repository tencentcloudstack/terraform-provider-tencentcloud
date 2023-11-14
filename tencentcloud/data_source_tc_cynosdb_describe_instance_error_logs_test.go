package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbDescribeInstanceErrorLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbDescribeInstanceErrorLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_describe_instance_error_logs.describe_instance_error_logs")),
			},
		},
	})
}

const testAccCynosdbDescribeInstanceErrorLogsDataSource = `

data "tencentcloud_cynosdb_describe_instance_error_logs" "describe_instance_error_logs" {
  instance_id = "cynosdbmysql-ins-4senc2fm"
  start_time = "2022-01-02 15:04:05"
  end_time = "2022-02-02 15:04:05"
  order_by = "Timestamp"
  order_by_type = "ASC"
  log_levels = 
  key_words = 
  }

`
