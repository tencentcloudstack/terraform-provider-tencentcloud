package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudClbListenersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListenersDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_listeners.listeners", "listener_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_listeners.listeners", "listener_list.0.listener_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.0.listener_name", "mylistener1234"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.0.port", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.0.session_expire_time", "30"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_listeners.listeners", "listener_list.0.scheduler", "WRR"),
				),
			},
		},
	})
}

const testAccClbListenersDataSource = `
resource "tencentcloud_clb_instance" "clb" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener" {
  clb_id              = "${tencentcloud_clb_instance.clb.id}"
  port                = 1
  protocol            = "TCP"
  listener_name       = "mylistener1234"
  session_expire_time = 30
  scheduler           = "WRR"
}

data "tencentcloud_clb_listeners" "listeners" {
  clb_id      = "${tencentcloud_clb_instance.clb.id}"
  listener_id = "${tencentcloud_clb_listener.listener.id}"
}
`
