package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmExportImagesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmExportImages,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "bucket_name"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "image_id", "img-l7uxaine"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "file_name_prefix", "test-"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "export_format", "RAW"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "only_export_root_disk", "false"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "dry_run", "false"),
				),
			},
		},
	})
}

const testAccCvmExportImagesBasis = `
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "private_sbucket" {
  bucket      = "tf-private-bucket-${local.app_id}"
  acl         = "private"
  force_clean = true
}
`

const testAccCvmExportImages = testAccCvmExportImagesBasis + `
resource "tencentcloud_cvm_export_images" "export_images" {
  bucket_name           = tencentcloud_cos_bucket.private_sbucket.bucket
  image_id              = "img-l7uxaine"
  file_name_prefix      = "test-"
  export_format         = "RAW"
  only_export_root_disk = false
  dry_run               = false
}
`
