package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudMonitorScrapeJob_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScrapeJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testScrapeJob_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScrapeJobExists("tencentcloud_monitor_tmp_scrape_job.basic"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_scrape_job.basic", "config", "job_name: demo-config-test\nhonor_timestamps: true\nmetrics_path: /metrics\nscheme: https\n"),
				),
			},
			{
				Config: testScrapeJob_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScrapeJobExists("tencentcloud_monitor_tmp_scrape_job.update"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_scrape_job.update", "config", "job_name: demo-config-test-update\nhonor_timestamps: true\nmetrics_path: /metrics\nscheme: https\n"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_scrape_job.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckScrapeJobDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MonitorService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_scrape_job" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		tmpScrapeJob, err := service.DescribeMonitorTmpScrapeJob(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if tmpScrapeJob != nil {
			return fmt.Errorf("scrape job %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckScrapeJobExists(r string) resource.TestCheckFunc {
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
		tmpScrapeJob, err := service.DescribeMonitorTmpScrapeJob(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if tmpScrapeJob == nil {
			return fmt.Errorf("scrape job %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testScrapeJobVar = `
variable "prometheus_id" {
  default = "` + defaultPrometheusId + `"
}
variable "agent_id" {
  default = "` + defaultAgentId + `"
}
`
const testScrapeJob_basic = testScrapeJobVar + `
resource "tencentcloud_monitor_tmp_scrape_job" "basic" {
  instance_id 	= var.prometheus_id
  agent_id 		= var.agent_id
  config 		= <<-EOT
job_name: demo-config-test
honor_timestamps: true
metrics_path: /metrics
scheme: https
EOT
}`

const testScrapeJob_update = testScrapeJobVar + `
resource "tencentcloud_monitor_tmp_scrape_job" "update" {
  instance_id	= var.prometheus_id
  agent_id 		= var.agent_id
  config 		= <<-EOT
job_name: demo-config-test-update
honor_timestamps: true
metrics_path: /metrics
scheme: https
EOT
}`
