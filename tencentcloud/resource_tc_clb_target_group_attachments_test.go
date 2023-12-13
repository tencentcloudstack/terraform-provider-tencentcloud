package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClbTargetGroupAttachmentsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbTargetGroupAttachments,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "load_balancer_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_target_group_attachments.target_group_attachments", "associations.#"),
				),
			},
		},
	})
}

const testAccClbTargetGroupAttachments = `

resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf_test_clb_attach"
  vpc_id = "vpc-5kwngvex"
}

resource "tencentcloud_clb_listener" "public_listeners" {
  clb_id        = tencentcloud_clb_instance.clb_basic.id
  #  protocol      = "HTTPS"
  #  port          = "443"
  protocol      = "HTTP"
  port          = "8090"
  listener_name = "iac-test-attach-2"
}

resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic2" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "baidu.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic3" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "tencent.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_listener_rule" "rule_basic4" {
  clb_id              = tencentcloud_clb_instance.clb_basic.id
  listener_id         = tencentcloud_clb_listener.public_listeners.listener_id
  domain              = "aws.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
  target_type         = "TARGETGROUP"
}
resource "tencentcloud_clb_target_group_attachments" "target_group_attachments" {
  load_balancer_id = tencentcloud_clb_instance.clb_basic.id
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-ln4gk8me"
    location_id = tencentcloud_clb_listener_rule.rule_basic.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-iv3mdrfy"
    location_id = tencentcloud_clb_listener_rule.rule_basic2.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-ctz951k0"
    location_id = tencentcloud_clb_listener_rule.rule_basic3.rule_id
  }
  associations  {
    listener_id = tencentcloud_clb_listener.public_listeners.listener_id
    target_group_id = "lbtg-kukb1j9c"
    location_id = tencentcloud_clb_listener_rule.rule_basic4.rule_id
  }
  depends_on = [tencentcloud_clb_listener.public_listeners,
    tencentcloud_clb_listener_rule.rule_basic4,
    tencentcloud_clb_listener_rule.rule_basic3,
    tencentcloud_clb_listener_rule.rule_basic2,
    tencentcloud_clb_listener_rule.rule_basic]
}

`
