package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmImportImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImportImage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_import_image.import_image", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_import_image.import_image",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmImportImage = `

resource "tencentcloud_cvm_import_image" "import_image" {
  architecture = "x86_64"
  os_type = "CentOS"
  os_version = "7"
  image_url = "http://111-1251233127.cosd.myqcloud.com/Windows%20Server%202008%20R2%20x64a.vmdk"
  image_name = "sample"
  image_description = "sampleimage"
  dry_run = false
  force = false
  tag_specification {
		resource_type = "image"
		tags {
			key = "tagKey"
			value = "tagValue"
		}

  }
  license_type = "TencentCloud"
  boot_mode = "Legacy BIOS"
  tags = {
    "createdBy" = "terraform"
  }
}

`
