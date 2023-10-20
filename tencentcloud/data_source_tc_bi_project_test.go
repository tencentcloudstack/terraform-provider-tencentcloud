package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudBiProjectDataSource_basic -v
func TestAccTencentCloudBiProjectDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBiProjectDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_bi_project.project"),
					resource.TestCheckResourceAttr("data.tencentcloud_bi_project.project", "list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.apply"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.auth_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.color_code"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.corp_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.created_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.created_user"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.is_external_manage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.member_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.page_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.seed"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_bi_project.project", "list.0.updated_user"),
				),
			},
		},
	})
}

const testAccBiProjectDataSource = `

data "tencentcloud_bi_project" "project" {
  keyword = "keep"
  all_page = true
}

`
