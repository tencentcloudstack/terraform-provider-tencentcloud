package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudPtsTmpKeyResource_basic -v
func TestAccTencentCloudPtsTmpKeyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsTmpKey,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "scenario_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "expired_time"),
					resource.TestCheckResourceAttr("tencentcloud_pts_tmp_key_generate.tmp_key", "credentials.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "credentials.0.tmp_secret_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "credentials.0.tmp_secret_key"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_tmp_key_generate.tmp_key", "credentials.0.token"),
				),
			},
		},
	})
}

const testAccPtsTmpKeyVar = `
variable "project_id" {
  default = "` + defaultPtsProjectId + `"
}
variable "scenario_id" {
	default = "` + defaultScenarioId + `"
}
  
`

const testAccPtsTmpKey = testAccPtsTmpKeyVar + `
resource "tencentcloud_pts_tmp_key_generate" "tmp_key" {
  project_id  = var.project_id
  scenario_id = var.scenario_id
}
`
