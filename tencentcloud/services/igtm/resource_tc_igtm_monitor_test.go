package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmMonitorResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIgtmMonitor,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "monitor_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "check_protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "check_interval"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "timeout"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "fail_times"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "fail_rate"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "detector_style"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "detector_group_ids"),
				),
			},
			{
				Config: testAccIgtmMonitorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "monitor_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "check_protocol"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "check_interval"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "timeout"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "fail_times"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "fail_rate"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "detector_style"),
					resource.TestCheckResourceAttrSet("tencentcloud_igtm_monitor.example", "detector_group_ids"),
				),
			},
			{
				ResourceName:      "tencentcloud_igtm_monitor.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIgtmMonitor = `
resource "tencentcloud_igtm_monitor" "example" {
  monitor_name          = "tf-example"
  check_protocol        = "PING"
  check_interval        = 60
  timeout               = 5
  fail_times            = 1
  fail_rate             = 50
  detector_style        = "INTERNAL"
  detector_group_ids    = [30, 31, 32, 34, 37, 38, 39, 1, 2, 3, 7, 8, 9, 10, 11, 12]
  ping_num              = 20
  tcp_port              = 443
  path                  = "/"
  return_code_threshold = 500
  enable_redirect       = "DISABLED"
  enable_sni            = "DISABLED"
  packet_loss_rate      = 90
  continue_period       = 1
}
`

const testAccIgtmMonitorUpdate = `
resource "tencentcloud_igtm_monitor" "example" {
  monitor_name          = "tf-example-update"
  check_protocol        = "PING"
  check_interval        = 60
  timeout               = 10
  fail_times            = 1
  fail_rate             = 50
  detector_style        = "INTERNAL"
  detector_group_ids    = [30, 31, 32, 34, 37, 38, 39, 1, 2, 3, 7, 8, 9, 10, 11, 12]
  ping_num              = 20
  tcp_port              = 443
  path                  = "/"
  return_code_threshold = 500
  enable_redirect       = "DISABLED"
  enable_sni            = "DISABLED"
  packet_loss_rate      = 90
  continue_period       = 1
}
`
