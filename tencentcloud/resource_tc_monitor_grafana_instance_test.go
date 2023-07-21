package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorGrafanaInstance_basic -v
func TestAccTencentCloudMonitorGrafanaInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGrafanaInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaInstanceExists("tencentcloud_monitor_grafana_instance.grafanaInstance"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "instance_name", "test-grafana"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "vpc_id", defaultGrafanaVpcId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "enable_internet", "false"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_grafana_instance.grafanaInstance", "internal_url"),
				),
			},
			{
				Config: testAccMonitorGrafanaInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaInstanceExists("tencentcloud_monitor_grafana_instance.grafanaInstance"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "instance_name", "test-grafana-update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "vpc_id", defaultGrafanaVpcId),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_instance.grafanaInstance", "enable_internet", "false"),
				),
			},
			{
				ResourceName:            "tencentcloud_monitor_grafana_instance.grafanaInstance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grafana_init_password", "is_destroy"},
			},
		},
	})
}

func testAccCheckGrafanaInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_grafana_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		instanceId := rs.Primary.ID

		instance, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
		if err != nil {
			return err
		}

		if instance == nil {
			return nil
		}

		if instance != nil && *instance.InstanceStatus < 5 {
			return fmt.Errorf("GrafanaInstance %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckGrafanaInstanceExists(r string) resource.TestCheckFunc {
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
		instance, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
		if err != nil {
			return err
		}

		if instance == nil {
			return fmt.Errorf("GrafanaInstance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testMonitorGrafanaInstanceVar = `
variable "vpc_id" {
  default = "` + defaultGrafanaVpcId + `"
}
variable "subnet_id" {
  default = "` + defaultGrafanaSubnetId + `"
}
`

const testAccMonitorGrafanaInstance = testMonitorGrafanaInstanceVar + `

resource "tencentcloud_monitor_grafana_instance" "grafanaInstance" {
  instance_name = "test-grafana"
  vpc_id = var.vpc_id
  subnet_ids = [var.subnet_id]
  grafana_init_password = "1234567890"
  enable_internet = false
  is_destroy = true

  tags = {
    "createdBy" = "test"
  }
}
`
const testAccMonitorGrafanaInstance_update = testMonitorGrafanaInstanceVar + `

resource "tencentcloud_monitor_grafana_instance" "grafanaInstance" {
  instance_name = "test-grafana-update"
  vpc_id = var.vpc_id
  subnet_ids = [var.subnet_id]
  grafana_init_password = "1234567890"
  enable_internet = false
  is_destroy = true

  tags = {
    "createdBy" = "test"
  }
}
`
