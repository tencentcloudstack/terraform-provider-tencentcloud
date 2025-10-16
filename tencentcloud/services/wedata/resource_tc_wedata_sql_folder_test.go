package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSqlFolderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataSqlFolder,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "folder_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "access_scope"),
				),
			},
			{
				Config: testAccWedataSqlFolderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "folder_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_sql_folder.example", "access_scope"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_sql_folder.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataSqlFolder = `
resource "tencentcloud_wedata_sql_folder" "example" {
  folder_name        = "tf_example"
  project_id         = "2983848457986924544"
  parent_folder_path = "/"
  access_scope       = "SHARED"
}
`

const testAccWedataSqlFolderUpdate = `
resource "tencentcloud_wedata_sql_folder" "example" {
  folder_name        = "tf_example_update"
  project_id         = "2983848457986924544"
  parent_folder_path = "/"
  access_scope       = "PRIVATE"
}
`
