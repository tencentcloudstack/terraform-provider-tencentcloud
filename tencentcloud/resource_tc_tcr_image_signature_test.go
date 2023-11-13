package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrImageSignatureResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImageSignature,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_image_signature.image_signature", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_image_signature.image_signature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrImageSignature = `

resource "tencentcloud_tcr_image_signature" "image_signature" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
  tags = {
    "createdBy" = "terraform"
  }
}

`
