package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataCodeFolderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataCodeFolder,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "folder_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "parent_folder_path"),
				),
			},
			{
				Config: testAccWedataCodeFolderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "folder_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_folder.example", "parent_folder_path"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_code_folder.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataCodeFolder = `
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example"
  parent_folder_path = "/"
}
`

const testAccWedataCodeFolderUpdate = `
resource "tencentcloud_wedata_code_folder" "example" {
  project_id         = "2983848457986924544"
  folder_name        = "tf_example_update"
  parent_folder_path = "/"
}
`
