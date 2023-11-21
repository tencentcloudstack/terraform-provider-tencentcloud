package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTestingClbListener_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTestingListener_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "listener_name", "listener_basic"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "scheduler", "WRR"),
				),
			},
			{
				Config: testAccClbTestingListenerUpdate_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "listener_name", "listener_basic_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "port", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "scheduler", "WRR"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_listener.listener_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbTestingListener_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  protocol            = "TCP"
  listener_name       = "listener_basic"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
`
const testAccClbTestingListenerUpdate_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-listener-basic"
}
resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  port                = 1
  protocol            = "TCP"
  listener_name       = "listener_basic_update"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
`
