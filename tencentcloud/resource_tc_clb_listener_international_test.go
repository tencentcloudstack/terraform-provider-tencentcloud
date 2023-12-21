package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudInternationalClbListener_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInternationalClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInternationalClbListener_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInternationalClbListenerExists("tencentcloud_clb_listener.listener_basic"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_basic", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_basic", "listener_name", "listener_basic"),
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

func testAccCheckInternationalClbListenerDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_listener" {
			continue
		}
		time.Sleep(5 * time.Second)
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, FILED_SP)
		itemLength := len(items)
		listenerId := items[itemLength-1]
		clbId := rs.Primary.Attributes["clb_id"]
		if itemLength == 2 && clbId != "" {
			clbId = items[0]
		}
		instance, err := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if instance != nil && err == nil {
			return fmt.Errorf("[CHECK][CLB listener][Destroy] check: CLB listener still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckInternationalClbListenerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLB listener][Exists] check: CLB listener %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLB listener][Exists] check: CLB listener id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		resourceId := rs.Primary.ID
		items := strings.Split(resourceId, FILED_SP)
		itemLength := len(items)
		listenerId := items[itemLength-1]
		clbId := rs.Primary.Attributes["clb_id"]
		if itemLength == 2 && clbId != "" {
			clbId = items[0]
		}
		instance, err := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLB listener][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccInternationalClbListener_basic = `
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
