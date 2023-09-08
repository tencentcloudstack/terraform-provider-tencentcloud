package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudRumPerformancePageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumPerformancePageDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_performance_page.performance_page")),
			},
		},
	})
}

const testAccRumPerformancePageDataSource = `

data "tencentcloud_rum_performance_page" "performance_page" {
  i_d = 1
  start_time = 1625444040
  end_time = 1625454840
  type = "pagepv"
  level = "1"
  isp = "中国电信"
  area = "广州市"
  net_type = "2"
  platform = "2"
  device = "Apple - iPhone"
  version_num = "1.0"
  ext_first = "ext1"
  ext_second = "ext2"
  ext_third = "ext3"
  is_abroad = "0"
  browser = "Chrome(79.0)"
  os = "Windows - 10"
  engine = "Blink(79.0)"
  brand = "Apple"
  from = "https://user.qzone.qq.com/"
  cost_type = "50"
  env = "production"
  net_status = "0"
  }

`
