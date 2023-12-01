package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataScriptResource_basic -v
func TestAccTencentCloudNeedFixWedataScriptResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataScript,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "file_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "bucket_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "region"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "file_extension_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_script.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataScriptUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "file_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "bucket_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "region"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_script.example", "file_extension_type"),
				),
			},
		},
	})
}

const testAccWedataScript = `
resource "tencentcloud_wedata_script" "example" {
  file_path           = "/datastudio/project/tf_example.sql"
  project_id          = "1470575647377821696"
  bucket_name         = "wedata-demo-1257305158"
  region              = "ap-guangzhou"
  file_extension_type = "sql"
}
`

const testAccWedataScriptUpdate = `
resource "tencentcloud_wedata_script" "example" {
  file_path           = "/datastudio/project/tf_example.sql"
  project_id          = "1470575647377821696"
  bucket_name         = "wedata-demo-1257305158"
  region              = "ap-guangzhou"
  file_extension_type = "sql"
}
`
