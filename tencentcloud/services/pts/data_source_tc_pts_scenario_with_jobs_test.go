package pts_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudPtsScenarioWithJobsDataSource_basic -v
func TestAccTencentCloudPtsScenarioWithJobsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsScenarioWithJobsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.created_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.geo_regions_load_distribution.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.geo_regions_load_distribution.0.percentage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.geo_regions_load_distribution.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.geo_regions_load_distribution.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.0.concurrency.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.0.concurrency.0.graceful_stop_seconds"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.0.concurrency.0.stages.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.0.concurrency.0.stages.0.duration_seconds"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.load.0.load_spec.0.concurrency.0.stages.0.target_virtual_users"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.project_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.scenario_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.sub_account_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_pts_scenario_with_jobs.scenario_with_jobs", "scenario_with_jobs_set.0.scenario.0.updated_at"),
				),
			},
		},
	})
}

const testAccPtsScenarioWithJobsDataSource = `

data "tencentcloud_pts_scenario_with_jobs" "scenario_with_jobs" {
  project_ids    = ["project-45vw7v82"]
  scenario_ids   = ["scenario-koakp3h6"]
  scenario_name  = "pts-jmeter"
  ascend         = true
  ignore_script  = true
  ignore_dataset = true
  scenario_type  = "pts-jmeter"
}

`
