package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbProjectsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbProjectsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_projects.projects"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.owner_uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.src_app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_projects.projects", "projects.0.status"),
				),
			},
		},
	})
}

const testAccDcdbProjectsDataSource = `

data "tencentcloud_dcdb_projects" "projects" {}

`
