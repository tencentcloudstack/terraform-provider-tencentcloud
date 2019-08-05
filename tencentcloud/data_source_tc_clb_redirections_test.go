package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudClbRedirectionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbRedirectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbRedirectionsDataSource,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbRedirectionExists("tencentcloud_clb_redirection.redirection_basic"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_redirections.redirections", "redirection_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_redirections.redirections", "redirection_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_redirections.redirections", "redirection_list.0.source_listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_redirections.redirections", "redirection_list.0.target_listener_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_redirections.redirections", "redirection_list.0.rewrite_source_loc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_redirections.redirections", "redirection_list.0.rewrite_target_loc_id"),
				),
			},
		},
	})
}

const testAccClbRedirectionsDataSource = `
resource "tencentcloud_clb_instance" "clb" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = "${tencentcloud_clb_instance.clb.id}"
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}


resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}


resource "tencentcloud_clb_listener" "listener_target" {
  clb_id        = "${tencentcloud_clb_instance.clb.id}"
  port          = 44
  protocol      = "HTTP"
  listener_name = "listener_basic1"
}


resource "tencentcloud_clb_listener_rule" "rule_target" {
  clb_id              = "${tencentcloud_clb_instance.clb.id}"
  listener_id         = "${tencentcloud_clb_listener.listener_target.id}"
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
resource "tencentcloud_clb_redirection" "redirection_basic" {
  clb_id                = "${tencentcloud_clb_instance.clb.id}"
  source_listener_id    = "${tencentcloud_clb_listener.listener_basic.id}"
  target_listener_id    = "${tencentcloud_clb_listener.listener_target.id}"
  rewrite_source_loc_id = "${tencentcloud_clb_listener_rule.rule_basic.id}"
  rewrite_target_loc_id = "${tencentcloud_clb_listener_rule.rule_target.id}"
}

data "tencentcloud_clb_redirections" "redirections" {
  clb_id                = "${tencentcloud_clb_instance.clb.id}"
  source_listener_id    = "${tencentcloud_clb_redirection.redirection_basic.source_listener_id}"
  rewrite_source_loc_id = "${tencentcloud_clb_redirection.redirection_basic.rewrite_source_loc_id}"
}

`
