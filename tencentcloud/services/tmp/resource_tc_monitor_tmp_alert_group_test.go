package tmp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -test.run TestAccTencentCloudMonitorTmpAlertGroupResource_basic -v
func TestAccTencentCloudMonitorTmpAlertGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckAlertGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpAlertGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertGroupExists("tencentcloud_monitor_tmp_alert_group.tmp_alert_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_alert_group.tmp_alert_group", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_alert_group.tmp_alert_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitorTmpAlertGroupUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlertGroupExists("tencentcloud_monitor_tmp_alert_group.tmp_alert_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_alert_group.tmp_alert_group", "id"),
				),
			},
		},
	})
}

func testAccCheckAlertGroupDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_alert_group" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		group, err := service.DescribeMonitorTmpAlertGroupById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if group != nil {
			return fmt.Errorf("group %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAlertGroupExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		group, err := service.DescribeMonitorTmpAlertGroupById(ctx, ids[0], ids[1])
		if err != nil {
			return err
		}

		if group == nil {
			return fmt.Errorf("group %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMonitorTmpAlertGroup = testInstance_basic + `

resource "tencentcloud_monitor_tmp_alert_group" "tmp_alert_group" {
  amp_receivers = [
    "notice-om017kc2",
  ]
  group_name      = "tf-test"
  instance_id     = tencentcloud_monitor_tmp_instance.example.id
  repeat_interval = "5m"

  custom_receiver {
    type = "amp"
  }

  rules {
    duration  = "1m"
    expr      = "up{job=\"prometheus-agent\"} != 1"
    rule_name = "Agent health check"
    state     = 2

    annotations = {
      "summary"     = "Agent health check"
      "description" = "Agent {{$labels.instance}} is deactivated, please pay attention!"
    }

    labels = {
      "severity" = "critical"
    }
  }
}

`

const testAccMonitorTmpAlertGroupUp = testInstance_basic + `

resource "tencentcloud_monitor_tmp_alert_group" "tmp_alert_group" {
  amp_receivers = [
    "notice-om017kc2",
  ]
  group_name      = "tf-test-up"
  instance_id     = tencentcloud_monitor_tmp_instance.example.id
  repeat_interval = "1h"

  custom_receiver {
    type = "amp"
  }

  rules {
    duration  = "1m"
    expr      = "up{job=\"prometheus-agent\"} != 1"
    rule_name = "Agent health check up"
    state     = 2

    annotations = {
      "summary"     = "Agent health check"
      "description" = "Agent {{$labels.instance}} is deactivated, please pay attention!!"
    }

    labels = {
      "severity" = "critical"
    }
  }
}

`
