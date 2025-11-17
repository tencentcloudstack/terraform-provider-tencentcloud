package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafOwaspWhiteRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafOwaspWhiteRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_white_rule.example", "id"),
				),
			},
			{
				Config: testAccWafOwaspWhiteRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_owasp_white_rule.example", "id"),
				),
			},
		},
	})
}

const testAccWafOwaspWhiteRule = `
resource "tencentcloud_waf_owasp_white_rule" "example" {
  name   = "tf-example"
  domain = "example.qcloud.com"
  strategies {
    field              = "IP"
    compare_func       = "ipmatch"
    content            = "1.1.1.1"
    arg                = ""
    case_not_sensitive = 0
  }
  ids = [
    10000000,
    20000000,
    30000000,
    40000000,
    90000000,
    110000000,
    190000000,
    200000000,
    210000000,
    220000000,
    230000000,
    240000000,
    250000000,
    260000000,
    270000000,
    280000000,
    290000000,
    300000000,
    310000000,
    320000000,
    330000000,
    340000000,
    350000000,
    360000000,
    370000000
  ]
  type     = 1
  job_type = "TimedJob"
  job_date_time {
    timed {
      start_date_time = 0
      end_date_time   = 0
    }

    time_t_zone = "UTC+8"
  }
  expire_time = 0
  status      = 1
}
`

const testAccWafOwaspWhiteRuleUpdate = `
resource "tencentcloud_waf_owasp_white_rule" "example" {
  name   = "tf-example"
  domain = "example.qcloud.com"
  strategies {
    field              = "IP"
    compare_func       = "ipmatch"
    content            = "1.1.1.1"
    arg                = ""
    case_not_sensitive = 0
  }
  ids = [
    10000000,
    20000000,
    30000000,
    40000000,
    90000000,
    110000000,
    190000000,
    200000000,
  ]
  type     = 1
  job_type = "TimedJob"
  job_date_time {
    timed {
      start_date_time = 0
      end_date_time   = 0
    }

    time_t_zone = "UTC+8"
  }
  expire_time = 0
  status      = 0
}
`
