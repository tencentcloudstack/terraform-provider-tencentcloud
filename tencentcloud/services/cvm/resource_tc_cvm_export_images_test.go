package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmExportImagesResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmExportImagesResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "only_export_root_disk", "false"), resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "dry_run", "false"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_export_images.export_images", "bucket_name"), resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "image_id", "img-l7uxaine"), resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "file_name_prefix", "test-"), resource.TestCheckResourceAttr("tencentcloud_cvm_export_images.export_images", "export_format", "RAW")),
			},
		},
	})
}

const testAccCvmExportImagesResource_BasicCreate = `

data "tencentcloud_user_info" "info" {
}
resource "tencentcloud_cos_bucket" "private_sbucket" {
    bucket = "tf-private-bucket-${data.tencentcloud_user_info.info.app_id}"
    acl = "private"
    force_clean = true
}
resource "tencentcloud_cvm_export_images" "export_images" {
    bucket_name = tencentcloud_cos_bucket.private_sbucket.bucket
    image_id = "img-l7uxaine"
    file_name_prefix = "test-"
    export_format = "RAW"
    only_export_root_disk = false
    dry_run = false
}

`
