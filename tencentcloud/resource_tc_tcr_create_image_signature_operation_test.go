package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: fmt.Sprintf(testAccTcrImageSignatureOperation, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_create_image_signature_operation.sign_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_create_image_signature_operation.sign_operation", "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_create_image_signature_operation.sign_operation", "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr("tencentcloud_tcr_create_image_signature_operation.sign_operation", "repository_name", defaultTCRRepoName),
					resource.TestCheckResourceAttr("tencentcloud_tcr_create_image_signature_operation.sign_operation", "image_version", "v1"),
				),
			},
		},
	})
}

const testAccTcrImageSignatureOperation = `

resource "tencentcloud_tcr_create_image_signature_operation" "sign_operation" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  image_version = "v1"
}

`
