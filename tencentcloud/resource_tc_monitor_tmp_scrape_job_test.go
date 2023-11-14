package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMonitorTmpScrapeJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorTmpScrapeJob,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_scrape_job.tmp_scrape_job", "id")),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_scrape_job.tmp_scrape_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMonitorTmpScrapeJob = `

resource "tencentcloud_monitor_tmp_scrape_job" "tmp_scrape_job" {
  instance_id = "prom-dko9d0nu"
  agent_id = "agent-6a7g40k2"
  config = "job_name: demo-config"
}

`
