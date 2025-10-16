package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataResourceFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataResourceFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource_file.wedata_resource_file", "id")),
			},
			{
				Config: testAccWedataResourceFileUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_wedata_resource_file.wedata_resource_file", "resource_name", "tftest1.txt"),
				),
			},
		},
	})
}

const testAccWedataResourceFile = `
resource "tencentcloud_wedata_resource_file" "wedata_resource_file" {
  project_id         = 2905622749543821312
  resource_name      = "tftest.txt"
  bucket_name        = "data-manage-fsi-1315051789"
  cos_region         = "ap-beijing-fsi"
  parent_folder_path = "/"
  resource_file      = "/datastudio/resource/2905622749543821312/test"
}
`

const testAccWedataResourceFileUpdate = `
resource "tencentcloud_wedata_resource_file" "wedata_resource_file" {
  project_id         = 2905622749543821312
  resource_name      = "tftest1.txt"
  bucket_name        = "data-manage-fsi-1315051789"
  cos_region         = "ap-beijing-fsi"
  parent_folder_path = "/"
  resource_file      = "/datastudio/resource/2905622749543821312/test"
}
`
