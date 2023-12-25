package tmp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_record_rule_yaml
	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_record_rule_yaml", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_tke_record_rule_yaml",
		F:    testSweepRecordRule,
	})
}
func testSweepRecordRule(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(region)
	client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := svcmonitor.NewMonitorService(client)

	instanceId := tcacctest.DefaultPrometheusId

	response, err := service.DescribePrometheusRecordRuleByName(ctx, instanceId, "")
	if err != nil {
		return err
	}

	instances := response.Response.Records
	if len(instances) == 0 {
		//return fmt.Errorf("instance %s record rule not exist", recordRuleName)
		return nil
	}

	for _, record := range instances {
		err = service.DeletePrometheusRecordRuleYaml(ctx, instanceId, *record.Name)
		if err != nil {
			continue
		}
	}

	return nil
}

// go test -i; go test -test.run TestAccTencentCloudMonitorRecordRuleResource_basic -v
func TestAccTencentCloudMonitorRecordRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckRecordRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testRecordRuleYaml_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordRuleExists("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic", "name", "example-record-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic", "content", "apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: example-record-test\nspec:\n  groups:\n    - name: kube-apiserver.rules\n      rules:\n        - expr: sum(metrics_test)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d'"),
				),
			},
			{
				Config: testRecordRuleYaml_basic_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordRuleExists("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic", "name", "example-record-test"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_record_rule_yaml.basic", "content", "apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: example-record-test\nspec:\n  groups:\n    - name: kube-apiserver.rules\n      rules:\n        - expr: sum(metrics_test)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d'\n    - name: kube-apiserver.rules2\n      rules:\n        - expr: sum(metrics_test2)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d2'"),
				),
			},
			//{
			//	ResourceName:      "tencentcloud_monitor_tmp_tke_record_rule_yaml.foo",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

func testAccCheckRecordRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	recordService := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_record_rule_yaml" {
			continue
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		RecordRuleName := items[1]
		response, err := recordService.DescribePrometheusRecordRuleByName(ctx, instanceId, RecordRuleName)
		if len(response.Response.Records) > 0 {
			return fmt.Errorf("record rule %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckRecordRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("instance id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		recordRuleName := items[1]
		recordRuleService := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		response, err := recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, recordRuleName)
		if len(response.Response.Records) < 1 {
			return fmt.Errorf("record rule %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testRecordRuleYamlVar = `
variable "prometheus_id" {
  default = "` + tcacctest.DefaultPrometheusId + `"
}`

const testRecordRuleYaml_basic = testRecordRuleYamlVar + `
resource "tencentcloud_monitor_tmp_tke_record_rule_yaml" "basic" {
 instance_id       = var.prometheus_id
 content           = "apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: example-record-test\nspec:\n  groups:\n    - name: kube-apiserver.rules\n      rules:\n        - expr: sum(metrics_test)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d'"
}`

const testRecordRuleYaml_basic_update = testRecordRuleYamlVar + `
resource "tencentcloud_monitor_tmp_tke_record_rule_yaml" "basic" {
 instance_id       = var.prometheus_id
 content           = "apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: example-record-test\nspec:\n  groups:\n    - name: kube-apiserver.rules\n      rules:\n        - expr: sum(metrics_test)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d'\n    - name: kube-apiserver.rules2\n      rules:\n        - expr: sum(metrics_test2)\n          labels:\n            verb: read\n          record: 'apiserver_request:burnrate1d2'"
}`
