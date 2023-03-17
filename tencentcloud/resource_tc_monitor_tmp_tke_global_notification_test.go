package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorTmpTkeGlobalNotification_basic -v
func TestAccTencentCloudMonitorTmpTkeGlobalNotification_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTmpTkeGlobalNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTmpTkeGlobalNotification_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeGlobalNotificationExists("tencentcloud_monitor_tmp_tke_global_notification.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_global_notification.basic", "notification.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_global_notification.basic", "notification.0.type", "webhook"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_global_notification.basic", "notification.0.notify_way.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_global_notification.basic", "notification.0.phone_arrive_notice", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_tke_global_notification.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTmpTkeGlobalNotificationDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_global_notification" {
			continue
		}

		tmpGlobalNotification, err := service.DescribeTkeTmpGlobalNotification(ctx, rs.Primary.ID)
		if *tmpGlobalNotification.Enabled {
			return fmt.Errorf("global notification %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTmpTkeGlobalNotificationExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("instance id is not set")
		}

		tkeService := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		tmpGlobalNotification, err := tkeService.DescribeTkeTmpGlobalNotification(ctx, rs.Primary.ID)
		if !*tmpGlobalNotification.Enabled {
			return fmt.Errorf("global notification %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testTmpTkeGlobalNotificationVar = `
variable "prometheus_id" {
default = "` + defaultPrometheusId + `"
}
`

const testTmpTkeGlobalNotification_basic = testTmpTkeGlobalNotificationVar + `
resource "tencentcloud_monitor_tmp_tke_global_notification" "basic" {
 instance_id   = var.prometheus_id
 notification {
	enabled   	  		 = true
	type      	  		 = "webhook"
	alert_manager  {
     cluster_id   = ""
     cluster_type = ""
     url          = ""
	}
	web_hook			  = ""
	repeat_interval       = "5m"
	time_range_start      = "00:00:00"
	time_range_end        = "23:59:59"
	notify_way            = ["SMS", "EMAIL"]
	receiver_groups       = []
	phone_notify_order    = []
	phone_circle_times    = 0
	phone_inner_interval  = 0
	phone_circle_interval = 0
	phone_arrive_notice   = false
 }
}
`
