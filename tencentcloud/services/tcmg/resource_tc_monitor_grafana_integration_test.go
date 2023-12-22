package tcmg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorGrafanaIntegration_basic -v
func TestAccTencentCloudMonitorGrafanaIntegration_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckGrafanaIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaIntegration,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaIntegrationExists("tencentcloud_monitor_grafana_integration.grafanaIntegration"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_integration.grafanaIntegration", "instance_id", tcacctest.DefaultGrafanaInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_integration.grafanaIntegration", "kind", "tencentcloud-monitor-app"),
				),
			},
			{
				Config: testAccMonitorGrafanaIntegration_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaIntegrationExists("tencentcloud_monitor_grafana_integration.grafanaIntegration"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_integration.grafanaIntegration", "instance_id", tcacctest.DefaultGrafanaInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_integration.grafanaIntegration", "kind", "tencentcloud-monitor-app"),
				),
			},
		},
	})
}

func testAccCheckGrafanaIntegrationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_grafana_integration" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		integrationId := idSplit[0]
		instanceId := idSplit[1]

		instance, err := service.DescribeMonitorGrafanaIntegration(ctx, integrationId, instanceId)
		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("GrafanaIntegration %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckGrafanaIntegrationExists(r string) resource.TestCheckFunc {
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
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		integrationId := idSplit[0]
		instanceId := idSplit[1]

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorGrafanaIntegration(ctx, integrationId, instanceId)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("GrafanaIntegration %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testMonitorGrafanaIntegrationVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaIntegration = testMonitorGrafanaIntegrationVar + `

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  instance_id = var.instance_id
  kind 		  = "tencentcloud-monitor-app"
  content     = "{\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"authProvider\":{\"__anyOf\":\"使用密钥\",\"useRole\":true,\"secretId\":\"arunma@tencent.com\",\"secretKey\":\"123456789\"},\"name\":\"uint-test\"},\"grafanaSpec\":{\"organizationIds\":[]}}}"
}

`

const testAccMonitorGrafanaIntegration_update = testMonitorGrafanaIntegrationVar + `

resource "tencentcloud_monitor_grafana_integration" "grafanaIntegration" {
  content         = "{\"id\":\"integration-9st6kqz6\",\"kind\":\"tencentcloud-monitor-app\",\"spec\":{\"dataSourceSpec\":{\"name\":\"uint-test3\",\"authProvider\":{\"secretId\":\"ROLE\",\"useRole\":true,\"__anyOf\":\"使用服务角色\"}},\"grafanaSpec\":{\"organizationIds\":[]}}}"
  instance_id     = var.instance_id
  kind 		      = "tencentcloud-monitor-app"
}

`
