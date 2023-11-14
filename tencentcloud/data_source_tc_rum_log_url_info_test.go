package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumLogUrlInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogUrlInfoDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_url_info.log_url_info")),
			},
		},
	})
}

const testAccRumLogUrlInfoDataSource = `

data "tencentcloud_rum_log_url_info" "log_url_info" {
  i_d = 1
  start_time = 1625444040
  end_time = 1625454840
  }

`
