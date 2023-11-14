package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrDeleteImageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrDeleteImage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image.delete_image", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_delete_image.delete_image",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrDeleteImage = `

resource "tencentcloud_tcr_delete_image" "delete_image" {
  registry_id = "tcr-xxx"
  repository_name = "repo"
  image_version = "v1"
  namespace_name = "ns"
}

`
