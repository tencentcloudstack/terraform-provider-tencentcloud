package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrImageSignatureOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrImageSignatureOperation,
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.sign_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.sign_operation", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.sign_operation", "namespace_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.sign_operation", "repository_name"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_create_image_signature_operation.sign_operation", "image_version", "v1"),
				),
			},
		},
	})
}

const testAccTcrImageSignatureOperation = TCRDataSource + `

resource "tencentcloud_tcr_create_image_signature_operation" "sign_operation" {
  registry_id = local.tcr_id
  namespace_name = local.tcr_ns_name
  repository_name = local.tcr_repo
  image_version = "v1"
}

`
