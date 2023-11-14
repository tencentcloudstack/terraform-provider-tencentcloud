package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumWebVitalsPageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumWebVitalsPageDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_web_vitals_page.web_vitals_page")),
			},
		},
	})
}

const testAccRumWebVitalsPageDataSource = `

data "tencentcloud_rum_web_vitals_page" "web_vitals_page" {
  start_time = 1625444040
  end_time = 1625454840
  i_d = 1
  ext_second = "ext2"
  engine = "Blink(79.0)"
  isp = "中国电信"
  from = "https://user.qzone.qq.com/"
  level = "1"
  type = "from"
  brand = "Apple"
  area = "广州市"
  version_num = "1.0"
  platform = "2"
  ext_third = "ext3"
  ext_first = "ext1"
  net_type = "2"
  device = "Apple - iPhone"
  is_abroad = "0"
  os = "Windows - 10"
  browser = "Chrome(79.0)"
  cost_type = "50"
  env = "production"
  }

`
