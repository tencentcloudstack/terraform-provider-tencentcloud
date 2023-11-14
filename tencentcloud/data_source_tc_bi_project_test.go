package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_bi_project.project")),
			},
		},
	})
}

const testAccBiProjectDataSource = `

data "tencentcloud_bi_project" "project" {
  page_no = 1
  keyword = "abc"
  all_page = true
  module_collection = "sys_common_user"
      }

`
