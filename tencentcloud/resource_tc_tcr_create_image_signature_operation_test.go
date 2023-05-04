package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTcrImageSignatureOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImageSignatureOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.image_signature_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_create_image_signature_operation.image_signature_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrImageSignatureOperation = `

resource "tencentcloud_tcr_create_image_signature_operation" "image_signature_operation" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
}

`
