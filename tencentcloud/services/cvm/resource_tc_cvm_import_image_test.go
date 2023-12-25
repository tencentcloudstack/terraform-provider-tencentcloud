package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCvmImportImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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
  image_url = ""
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
