package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTcrDeleteImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrDeleteImageOperation,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "namespace_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "repository_name"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "image_version", "v2"),
				),
			},
		},
	})
}

const testAccTcrDeleteImageOperation = TCRDataSource + `

resource "tencentcloud_tcr_delete_image_operation" "delete_image_operation" {
  registry_id = local.tcr_id
  namespace_name = local.tcr_ns_name
  repository_name = local.tcr_repo
  image_version = "v2"
}

`
