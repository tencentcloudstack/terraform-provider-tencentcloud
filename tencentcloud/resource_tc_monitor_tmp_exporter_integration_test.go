package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMonitorExporterIntegration_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckExporterIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testExporterIntegration_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExporterIntegrationExists("tencentcloud_monitor_tmp_exporter_integration.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kind", "cvm-http-sd-exporter"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kube_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "cluster_id", "cls-ely08ic4"),
				),
			},
			{
				Config: testExporterIntegration_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExporterIntegrationExists("tencentcloud_monitor_tmp_exporter_integration.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kind", "cvm-http-sd-exporter"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kube_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "cluster_id", "cls-87o4klby"),
				),
			},
		},
	})
}

func testAccCheckExporterIntegrationDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_exporter_integration" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMonitorTmpExporterIntegration(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			return fmt.Errorf("ExporterIntegration %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckExporterIntegrationExists(r string) resource.TestCheckFunc {
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

		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		tmpExporterIntegration, err := service.DescribeMonitorTmpExporterIntegration(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if tmpExporterIntegration == nil {
			return fmt.Errorf("ExporterIntegration %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testExporterIntegrationVar = `
variable "prometheus_id" {
  default = "` + defaultPrometheusId + `"
}
variable "default_cluster" {
  default = "` + defaultTkeClusterId + `"
}
variable "cluster_id" {
  default = "` + tkeClusterIdAgent + `"
}
`
const testExporterIntegration_basic = testExporterIntegrationVar + `
resource "tencentcloud_monitor_tmp_exporter_integration" "basic" {
  instance_id	= var.prometheus_id	
  kind 			= "cvm-http-sd-exporter"
  content		= "{\"kind\":\"cvm-http-sd-exporter\",\"spec\":{\"job\":\"job_name: example-job-name\\nmetrics_path: /metrics\\ncvm_sd_configs:\\n- region: ap-guangzhou\\n  ports:\\n  - 9100\\n  filters:         \\n  - name: tag:示例标签键\\n    values: \\n    - 示例标签值\\nrelabel_configs: \\n- source_labels: [__meta_cvm_instance_state]\\n  regex: RUNNING\\n  action: keep\\n- regex: __meta_cvm_tag_(.*)\\n  replacement: $1\\n  action: labelmap\\n- source_labels: [__meta_cvm_region]\\n  target_label: region\\n  action: replace\"}}"
  kube_type		= 1
  cluster_id	= var.default_cluster
}`

const testExporterIntegration_update = testExporterIntegrationVar + `
resource "tencentcloud_monitor_tmp_exporter_integration" "basic" {
  instance_id	= var.prometheus_id	
  kind 			= "cvm-http-sd-exporter"
  content		= "{\"kind\":\"cvm-http-sd-exporter\",\"spec\":{\"job\":\"job_name: example-job-name\\nmetrics_path: /metrics\\ncvm_sd_configs:\\n- region: ap-guangzhou\\n  ports:\\n  - 9100\\n  filters:         \\n  - name: tag:示例标签键\\n    values: \\n    - 示例标签值\\nrelabel_configs: \\n- source_labels: [__meta_cvm_instance_state]\\n  regex: RUNNING\\n  action: keep\\n- regex: __meta_cvm_tag_(.*)\\n  replacement: $1\\n  action: labelmap\\n- source_labels: [__meta_cvm_region]\\n  target_label: region\\n  action: replace\"}}"
  kube_type		= 1
  cluster_id	= var.cluster_id
}`
