package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_monitor_tmp_tke_cluster_agent
	resource.AddTestSweepers("tencentcloud_monitor_tmp_tke_cluster_agent", &resource.Sweeper{
		Name: "tencentcloud_monitor_tmp_tke_cluster_agent",
		F:    testSweepClusterAgent,
	})
}
func testSweepClusterAgent(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(region)
	client := cli.(*TencentCloudClient).apiV3Conn
	service := MonitorService{client}

	instanceId := clusterPrometheusId
	clusterId := tkeClusterIdAgent
	clusterType := tkeClusterTypeAgent

	agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
	if err != nil {
		return err
	}

	if agents != nil {
		return nil
	}

	err = service.DeletePrometheusClusterAgent(ctx, instanceId, clusterId, clusterType)
	if err != nil {
		return err
	}

	return nil
}

// go test -i; go test -test.run TestAccTencentCloudMonitorClusterAgent_basic -v
func TestAccTencentCloudMonitorClusterAgent_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusterAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testClusterAgentYaml_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterAgentExists("tencentcloud_monitor_tmp_tke_cluster_agent.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_cluster_agent.basic", "agents.0.cluster_id", "cls-9ae9qo9k"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_tke_cluster_agent.basic", "agents.0.cluster_type", "eks"),
				),
			},
		},
	})
}

func testAccCheckClusterAgentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_tke_cluster_agent" {
			continue
		}
		items := strings.Split(rs.Primary.ID, FILED_SP)
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		clusterId := items[1]
		clusterType := items[2]
		agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if agents != nil {
			return fmt.Errorf("cluster agent %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckClusterAgentExists(r string) resource.TestCheckFunc {
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
		if len(items) != 3 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		clusterId := items[1]
		clusterType := items[2]
		service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if agents == nil {
			return fmt.Errorf("cluster agent %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testClusterAgentYamlVar = `
variable "prometheus_id" {
  default = "` + clusterPrometheusId + `"
}
variable "default_region" {
  default = "` + defaultRegion + `"
}
variable "agent_cluster_id" {
  default = "` + tkeClusterIdAgent + `"
}
variable "agent_cluster_type" {
  default = "` + tkeClusterTypeAgent + `"
}`

const testClusterAgentYaml_basic = testClusterAgentYamlVar + `
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "basic" {
  instance_id = var.prometheus_id
  agents {
    region          = var.default_region
    cluster_type    = var.agent_cluster_type
    cluster_id      = var.agent_cluster_id
    enable_external = false
  }
}`
