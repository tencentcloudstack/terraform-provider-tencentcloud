package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccTencentCloudClbListener_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_basic,
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
				ResourceName:      "tencentcloud_clb_listener.listener_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTencentCloudClblistener_tcp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_tcp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "30"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "100"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "2"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "2"),
				),
			},
			{
				Config: testAccClbListener_tcp_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_tcp"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_tcp", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "protocol", "TCP"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "listener_name", "listener_tcp_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "session_expire_time", "60"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "port", "44"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "scheduler", "WRR"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_switch", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_time_out", "20"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_interval_time", "200"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_health_num", "3"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_tcp", "health_check_unhealth_num", "3"),
				),
			},
		},
	})
}

func TestAccTencentCloudClblistener_https(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbListener_https,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "77"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_ssl_mode", "UNIDIRECTIONAL"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "certificate_id", "VfqcL1ME"),
				),
			},
			{
				Config: testAccClbListener_https_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClbListenerExists("tencentcloud_clb_listener.listener_https"),
					resource.TestCheckResourceAttrSet("tencentcloud_clb_listener.listener_https", "clb_id"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "protocol", "HTTPS"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "listener_name", "listener_https_update"),
					resource.TestCheckResourceAttr("tencentcloud_clb_listener.listener_https", "port", "33"),
				),
			},
		},
	})
}

func testAccCheckClbListenerDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbService := ClbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clb_listener" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, err := clbService.DescribeListenerById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb listener still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClbListenerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("clb listener %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("clb listener id is not set")
		}
		clbService := ClbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, err := clbService.DescribeListenerById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccClbListener_basic = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}

resource "tencentcloud_clb_listener" "listener_basic" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  port                = 1
  protocol            = "TCP"
  listener_name       = "listener_basic"
  session_expire_time = 30
  scheduler           = "WRR"
}
`

const testAccClbListener_tcp = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}
resource "tencentcloud_clb_listener" "listener_tcp" {
  clb_id                     = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name              = "listener_tcp"
  port                       = 44
  protocol                   = "TCP"
  health_check_switch        = 1
  health_check_time_out      = 30
  health_check_interval_time = 100
  health_check_health_num    = 2
  health_check_unhealth_num  = 2
  session_expire_time        = 30
  scheduler                  = "WRR"
}
`

const testAccClbListener_tcp_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}
resource "tencentcloud_clb_listener" "listener_tcp"{
        clb_id = "${tencentcloud_clb_instance.clb_basic.id}"
        listener_name              = "listener_tcp_update"
        port                       = 44
        protocol                   = "TCP"
        health_check_switch        = 1
        health_check_time_out      = 20
        health_check_interval_time = 200
        health_check_health_num    = 3
        health_check_unhealth_num  = 3
        session_expire_time        = 60
        scheduler                  = "WRR"
}
`

const testAccClbListener_https = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}
resource "tencentcloud_clb_listener" "listener_https" {
  clb_id               = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name        = "listener_https"
  port                 = 77
  protocol             = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VfqcL1ME"
}
`
const testAccClbListener_https_update = `
resource "tencentcloud_clb_instance" "clb_basic" {
  network_type = "OPEN"
  clb_name     = "tf-clb-basic"
}
resource "tencentcloud_clb_listener" "listener_https" {
  clb_id              = "${tencentcloud_clb_instance.clb_basic.id}"
  listener_name       = "listener_https_update"
  port                = 33
  protocol            = "HTTPS"
  certificate_ssl_mode = "UNIDIRECTIONAL"
  certificate_id       = "VfqcL1ME"
}
`
