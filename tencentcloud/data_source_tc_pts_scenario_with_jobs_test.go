package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPtsScenarioWithJobsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsScenarioWithJobsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs")),
			},
		},
	})
}

const testAccPtsScenarioWithJobsDataSource = `

data "tencentcloud_pts_scenario_with_jobs" "scenario_with_jobs" {
  project_ids = 
  scenario_ids = 
  scenario_name = "test"
  scenario_status = 2
  order_by = "updated_at"
  ascend = true
  scenario_related_jobs_params {
		offset = 2
		limit = 5
		order_by = "updated_at"
		ascend = true

  }
  ignore_script = true
  ignore_dataset = true
  scenario_type = "pts-http"
  owner = "tom"
  }

`
