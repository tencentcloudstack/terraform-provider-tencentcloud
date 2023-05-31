package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmExportImagesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmExportImages,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "id")),
			},
		},
	})
}

const testAccCvmExportImages = `
resource "tencentcloud_cvm_export_images" "export_images" {
	bucket_name = "keep-export-image-1308726196"
	image_id = "img-e4l9lc5o"
	file_name_prefix = "test-"
  }
`
