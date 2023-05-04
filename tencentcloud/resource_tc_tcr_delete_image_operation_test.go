package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTcrDeleteImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrDeleteImageOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_delete_image_operation.delete_image_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrDeleteImageOperation = `

resource "tencentcloud_tcr_delete_image_operation" "delete_image_operation" {
  registry_id = "tcr-xxx"
  repository_name = "repo"
  image_version = "v1"
  namespace_name = "ns"
}

`
