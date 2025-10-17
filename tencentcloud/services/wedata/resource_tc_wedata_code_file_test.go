package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataCodeFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataCodeFile,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_content"),
				),
			},
			{
				Config: testAccWedataCodeFileUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "parent_folder_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_file.example", "code_file_content"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_code_file.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataCodeFile = `
resource "tencentcloud_wedata_code_file" "example" {
  project_id         = ""
  code_file_name     = ""
  parent_folder_path = ""
  code_file_config {
    params = ""
    notebook_session_info {
      notebook_session_id   = ""
      notebook_session_name = ""
    }
  }

  code_file_content = ""
}
`

const testAccWedataCodeFileUpdate = `
resource "tencentcloud_wedata_code_file" "example" {
  project_id         = ""
  code_file_name     = ""
  parent_folder_path = ""
  code_file_config {
    params = ""
    notebook_session_info {
      notebook_session_id   = ""
      notebook_session_name = ""
    }
  }

  code_file_content = ""
}
`
