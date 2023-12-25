package tmp_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const DefaultKind = ""

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_exporter_integration
	resource.AddTestSweepers("tencentcloud_monitor_tmp_exporter_integration", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_exporter_integration",
		F:    testSweepExporterIntegration,
	})
}
func testSweepExporterIntegration(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcmonitor.NewMonitorService(client)

	instanceId := tcacctest.DefaultPrometheusId
	clusterId := tcacctest.TkeClusterIdAgent
	ids := strings.Join([]string{"", instanceId, strconv.Itoa(1), clusterId, DefaultKind}, tccommon.FILED_SP)

	for {
		instances, err := service.DescribeMonitorTmpExporterIntegration(ctx, ids)
		if err != nil {
			return err
		}

		if instances == nil {
			break
		}

		id := strings.Join([]string{*instances.Name, instanceId, strconv.Itoa(1), clusterId, DefaultKind}, tccommon.FILED_SP)
		err = service.DeleteMonitorTmpExporterIntegrationById(ctx, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestAccTencentCloudMonitorExporterIntegration_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckExporterIntegrationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testExporterIntegration_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExporterIntegrationExists("tencentcloud_monitor_tmp_exporter_integration.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kind", "cvm-http-sd-exporter"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kube_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "cluster_id", "cls-9ae9qo9k"),
				),
			},
			{
				Config: testExporterIntegration_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExporterIntegrationExists("tencentcloud_monitor_tmp_exporter_integration.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kind", "cvm-http-sd-exporter"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "kube_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_exporter_integration.basic", "cluster_id", "cls-9ae9qo9k"),
				),
			},
		},
	})
}

func testAccCheckExporterIntegrationDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
  default = "` + tcacctest.ClusterPrometheusId + `"
}
variable "cluster_id" {
  default = "` + tcacctest.TkeClusterIdAgent + `"
}
`
const testExporterIntegration_basic = testExporterIntegrationVar + `
resource "tencentcloud_monitor_tmp_exporter_integration" "basic" {
  instance_id	= var.prometheus_id	
  kind 			= "cvm-http-sd-exporter"
  content		= "{\"kind\":\"cvm-http-sd-exporter\",\"spec\":{\"job\":\"job_name: example-job-name-test\\nmetrics_path: /metrics\\ncvm_sd_configs:\\n- region: ap-guangzhou\\n  ports:\\n  - 9100\\n  filters:         \\n  - name: tag:示例标签键\\n    values: \\n    - 示例标签值\\nrelabel_configs: \\n- source_labels: [__meta_cvm_instance_state]\\n  regex: RUNNING\\n  action: keep\\n- regex: __meta_cvm_tag_(.*)\\n  replacement: $1\\n  action: labelmap\\n- source_labels: [__meta_cvm_region]\\n  target_label: region\\n  action: replace\"}}"
  kube_type		= 1
  cluster_id	= var.cluster_id
}`

const testExporterIntegration_update = testExporterIntegrationVar + `
resource "tencentcloud_monitor_tmp_exporter_integration" "basic" {
  instance_id	= var.prometheus_id	
  kind 			= "cvm-http-sd-exporter"
  content		= "{\"kind\":\"cvm-http-sd-exporter\",\"spec\":{\"job\":\"job_name: example-job-name-test\\nmetrics_path: /metrics\\ncvm_sd_configs:\\n- region: ap-guangzhou\\n  ports:\\n  - 9100\\n  filters:         \\n  - name: tag:示例标签键\\n    values: \\n    - 示例标签值\\nrelabel_configs: \\n- source_labels: [__meta_cvm_instance_state]\\n  regex: RUNNING\\n  action: keep\\n- regex: __meta_cvm_tag_(.*)\\n  replacement: $1\\n  action: labelmap\\n- source_labels: [__meta_cvm_region]\\n  target_label: region\\n  action: replace\"}}"
  kube_type		= 1
  cluster_id	= var.cluster_id
}`
