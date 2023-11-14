package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbErrorLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbErrorLogDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_error_log.error_log")),
			},
		},
	})
}

const testAccCdbErrorLogDataSource = `

data "tencentcloud_cdb_error_log" "error_log" {
  instance_id = ""
  start_time = 
  end_time = 
  key_words = 
  inst_type = ""
  }

`
