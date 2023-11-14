package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_list.log_list")),
			},
		},
	})
}

const testAccRumLogListDataSource = `

data "tencentcloud_rum_log_list" "log_list" {
  order_by = "desc"
  start_time = 1625444040000
  query = "id:123 AND type:&quot;log&quot;"
  end_time = 1625454840000
  i_d = 1
  }

`
