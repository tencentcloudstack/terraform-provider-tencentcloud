package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorTmpGrafanaConfigResource_basic -v
func TestAccTencentCloudMonitorTmpGrafanaConfigResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
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
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		instanceId := rs.Primary.ID

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
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
