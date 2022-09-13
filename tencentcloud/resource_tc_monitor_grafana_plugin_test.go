package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorGrafanaPlugin_basic -v
func TestAccTencentCloudMonitorGrafanaPlugin_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGrafanaPluginDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaPlugin,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaPluginExists("tencentcloud_monitor_grafana_plugin.grafanaPlugin"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_plugin.grafanaPlugin", "plugin_id", defaultGrafanaPlugin),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_plugin.grafanaPlugin", "version", defaultGrafanaVersion),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_plugin.grafanaPlugin",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGrafanaPluginDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_grafana_plugin" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		pluginId := idSplit[1]

		grafanaPlugin, err := service.DescribeMonitorGrafanaPlugin(ctx, instanceId, pluginId)
		if err != nil {
			return err
		}

		if grafanaPlugin != nil {
			return fmt.Errorf("GrafanaPlugin %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckGrafanaPluginExists(r string) resource.TestCheckFunc {
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
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		pluginId := idSplit[1]

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		grafanaPlugin, err := service.DescribeMonitorGrafanaPlugin(ctx, instanceId, pluginId)
		if err != nil {
			return err
		}

		if grafanaPlugin == nil {
			return fmt.Errorf("GrafanaPlugin %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testMonitorGrafanaPluginVar = `
variable "instance_id" {
  default = "` + defaultGrafanaInstanceId + `"
}
variable "plugin_id" {
  default = "` + defaultGrafanaPlugin + `"
}
variable "plugin_version" {
  default = "` + defaultGrafanaVersion + `"
}
`

const testAccMonitorGrafanaPlugin = testMonitorGrafanaPluginVar + `

resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
  version     = var.plugin_version
}

`
