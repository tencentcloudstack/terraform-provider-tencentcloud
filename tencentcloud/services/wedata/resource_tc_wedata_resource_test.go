package wedata_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixWedataResourceResource_basic -v
func TestAccTencentCloudNeedFixWedataResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "file_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "file_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "cos_bucket_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "cos_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "files_size"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_resource.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWedataResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "file_path"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "project_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "file_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "cos_bucket_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "cos_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_resource.example", "files_size"),
				),
			},
		},
	})
}

const testAccWedataResource = `
resource "tencentcloud_wedata_resource" "example" {
  file_path       = "/datastudio/resource/demo"
  project_id      = "1612982498218618880"
  file_name       = "tf_example"
  cos_bucket_name = "wedata-demo-1314991481"
  cos_region      = "ap-guangzhou"
  files_size      = "8165"
}
`

const testAccWedataResourceUpdate = `
resource "tencentcloud_wedata_resource" "example" {
  file_path       = "/datastudio/resource/demo"
  project_id      = "1612982498218618880"
  file_name       = "tf_example_update"
  cos_bucket_name = "wedata-demo-1314991481"
  cos_region      = "ap-guangzhou"
  files_size      = "7210"
}
`
