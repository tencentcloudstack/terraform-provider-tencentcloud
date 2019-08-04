package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudClbRewrite_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckClbRewriteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbRewrite_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbRewriteExists("tencentcloud_clb_rewrite.rewrite_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_rewrite.rewrite_basic", "clb_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_rewrite.rewrite_basic", "source_listener_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_rewrite.rewrite_basic", "target_listener_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_rewrite.rewrite_basic", "rewrite_source_loc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_rewrite.rewrite_basic", "rewrite_target_loc_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_clb_rewrite.rewrite_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClbRewriteDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_rewrite" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, err := clbService.DescribeRewriteInfoById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb rewrite still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbRewriteExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("clb rewrite %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("clb rewrite id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := clbService.DescribeRewriteInfoById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccClbRewrite_basic = `
resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id        = "lb-p7olt9e5"
  port          = 1
  protocol      = "HTTP"
  listener_name = "listener_basic"
}


resource "tencentcloud_clb_listener_rule" "rule_basic" {
  clb_id              = "lb-p7olt9e5"
  listener_id         = "${tencentcloud_clb_listener.listener_basic.id}"
  domain              = "abc.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}


resource "tencentcloud_clb_listener" "listener_target" {
  clb_id        = "lb-p7olt9e5"
  port          = 44
  protocol      = "HTTP"
  listener_name = "listener_basic1"
}


resource "tencentcloud_clb_listener_rule" "rule_target" {
  clb_id              = "lb-p7olt9e5"
  listener_id         = "${tencentcloud_clb_listener.listener_target.id}"
  domain              = "abcd.com"
  url                 = "/"
  session_expire_time = 30
  scheduler           = "WRR"
}
resource "tencentcloud_clb_rewrite" "rewrite_basic" {
  clb_id                = "lb-p7olt9e5"
  source_listener_id    = "${tencentcloud_clb_listener.listener_basic.id}"
  target_listener_id    = "${tencentcloud_clb_listener.listener_target.id}"
  rewrite_source_loc_id = "${tencentcloud_clb_listener_rule.rule_basic.id}"
  rewrite_target_loc_id = "${tencentcloud_clb_listener_rule.rule_target.id}"
}
`
