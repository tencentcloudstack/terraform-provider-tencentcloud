package tcmg_test

import (
	"context"
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorTmpGrafanaConfigResource_basic -v
func TestAccTencentCloudMonitorTmpGrafanaConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpGrafanaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorTmpGrafanaConfigExists("tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config", "config"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMonitorTmpGrafanaConfigExists(r string) resource.TestCheckFunc {
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
		instanceId := rs.Primary.ID

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorTmpGrafanaConfigById(ctx, instanceId)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("GrafanaInstance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMonitorTmpGrafanaConfig = testAccMonitorGrafanaInstance + `

resource "tencentcloud_monitor_tmp_grafana_config" "tmp_grafana_config" {
    config      = jsonencode(
        {
            server = {
                http_port           = 8080
                root_url            = "https://cloud-grafana.woa.com/grafana-ffrdnrfa/"
                serve_from_sub_path = true
            }
        }
    )
    instance_id = tencentcloud_monitor_grafana_instance.grafanaInstance.id
}

`
