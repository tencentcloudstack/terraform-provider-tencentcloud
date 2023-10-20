package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudBiUserProjectDataSource_basic -v
func TestAccTencentCloudBiUserProjectDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiUserProjectDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_bi_user_project.user_project"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.area_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.corp_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.created_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.created_user"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.first_modify"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.global_user_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.phone_number"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.updated_user"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.user_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_user_project.user_project", "list.0.user_name"),
				),
			},
		},
	})
}

const testAccBiUserProjectDataSource = `

data "tencentcloud_bi_user_project" "user_project" {
  project_id = 11015030
  all_page   = true
}

`
