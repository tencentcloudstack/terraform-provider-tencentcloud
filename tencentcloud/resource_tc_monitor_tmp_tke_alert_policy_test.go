package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//func init() {
//	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_alert_policy
//	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_alert_policy", &resource.Sweeper{
//		Name: "tencentcloud_monitor_tmp_tke_alert_policy",
//		F:    testSweepTmpTkeAlertPolicy,
//	})
//}
//
//func testSweepTmpTkeAlertPolicy(region string) error {
//	logId := getLogId(contextNil)
//	ctx := context.WithValue(context.TODO(), logIdKey, logId)
//	cli, _ := sharedClientForRegion(region)
//	client := cli.(*TencentCloudClient).apiV3Conn
//	service := TkeService{client}
//
//	instanceId := defaultPrometheusId
//
//	for {
//		tmpAlertPolicy, err := service.DescribeTkeTmpAlertPolicy(ctx, instanceId, "")
//		if err != nil {
//			return err
//		}
//
//		if tmpAlertPolicy == nil {
//			return nil
//		}
//
//		err = service.DeleteTkeTmpAlertPolicyById(ctx, instanceId, *tmpAlertPolicy.Id)
//		if err != nil {
//			return err
//		}
//	}
//}

func TestAccTencentCloudMonitorTmpTkeAlertPolicy_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTmpTkeAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTmpTkeAlertPolicy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmpTkeAlertPolicyExists("tencentcloud_monitor_tmp_tke_alert_policy.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.cluster_id", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.id", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.name", "alert_rule-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.enabled", "true"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.phone_arrive_notice", "false"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.phone_circle_interval", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.phone_circle_times", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.phone_inner_interval", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.repeat_interval", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.time_range_end", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.time_range_start", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.type", "amp"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.notification.0.web_hook", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.describe", ""),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.for", "5m"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.labels.0.name", "severity"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.labels.0.value", "warning"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.name", "rules-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.rule", "(count(kube_node_status_allocatable_cpu_cores) by (cluster) -1)   / count(kube_node_status_allocatable_cpu_cores) by (cluster)"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.rule_state", "0"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_alert_policy.basic", "alert_rule.0.rules.0.template", "集群{{ $labels.cluster }}内Pod申请的CPU过载，当前CPU申请占比{{ $value | humanizePercentage }}"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_tke_alert_policy.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTmpTkeAlertPolicyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_alert_policy" {
			continue
		}

		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tmpAlertPolicyId := items[1]

		tmpAlertPolicy, err := service.DescribeTkeTmpAlertPolicy(ctx, instanceId, tmpAlertPolicyId)
		if tmpAlertPolicy != nil {
			return fmt.Errorf("alert policy still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTmpTkeAlertPolicyExists(r string) resource.TestCheckFunc {
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
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tmpAlertPolicyId := items[1]

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		tmpAlertPolicy, err := service.DescribeTkeTmpAlertPolicy(ctx, instanceId, tmpAlertPolicyId)
		if tmpAlertPolicy == nil {
			return fmt.Errorf("alert policy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testTmpTkeAlertPolicyVar = `
variable "prometheus_id" {
  default = "` + defaultPrometheusId + `"
}
`
const testTmpTkeAlertPolicy_basic = testTmpTkeAlertPolicyVar + `
resource "tencentcloud_monitor_tmp_tke_alert_policy" "basic" {
  instance_id = var.prometheus_id
  alert_rule {
    name = "alert_rule-test"
    rules {
      name = "rules-test"
      rule = "(count(kube_node_status_allocatable_cpu_cores) by (cluster) -1)   / count(kube_node_status_allocatable_cpu_cores) by (cluster)"
      template = "集群{{ $labels.cluster }}内Pod申请的CPU过载，当前CPU申请占比{{ $value | humanizePercentage }}"
      for = "5m"
      labels {
        name  = "severity"
        value = "warning"
      }
    }
    notification {
      type = "amp"
      enabled = true
      alert_manager {
		url	= "xxx"
	  }
    }
  }
}`
