package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusFolderResource_basic -v
func TestAccTencentCloudNeedFixOceanusFolderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusFolder,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_folder.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "folder_name", "tf_example"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "parent_id", "folder-f40fq79g"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_folder.example", "folder_type"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "work_space_id", "space-bshmbms5"),
				),
			},
			{
				ResourceName:      "tencentcloud_oceanus_folder.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOceanusFolderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_folder.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "folder_name", "tf_example_update"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "parent_id", "folder-f40fq79g"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_folder.example", "folder_type"),
					resource.TestCheckResourceAttr("tencentcloud_oceanus_folder.example", "work_space_id", "space-bshmbms5"),
				),
			},
		},
	})
}

const testAccOceanusFolder = `
resource "tencentcloud_oceanus_folder" "example" {
  folder_name   = "tf_example"
  parent_id     = "folder-f40fq79g"
  folder_type   = 0
  work_space_id = "space-bshmbms5"
}
`

const testAccOceanusFolderUpdate = `
resource "tencentcloud_oceanus_folder" "example" {
  folder_name   = "tf_example_update"
  parent_id     = "folder-f40fq79g"
  folder_type   = 0
  work_space_id = "space-bshmbms5"
}
`
