package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmExportImagesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmExportImages,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_export_images.export_images",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmExportImages = `

resource "tencentcloud_cvm_export_images" "export_images" {
  bucket_name = "test-bucket-AppId"
  image_ids = 
  export_format = "RAW"
  file_name_prefix_list = 
  only_export_root_disk = true
  dry_run = false
  role_name = "CVM_QcsRole"
}

`
