package tencentcloud

import (
	"fmt"
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
				Config: fmt.Sprintf(testAccTcrDeleteImageOperation, defaultTCRInstanceId, defaultTCRNamespace, defaultTCRRepoName),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "registry_id", defaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "namespace_name", defaultTCRNamespace),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "repository_name", defaultTCRRepoName),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "image_version", "v2"),
				),
			},
		},
	})
}

const testAccTcrDeleteImageOperation = `

resource "tencentcloud_tcr_delete_image_operation" "delete_image_operation" {
  registry_id = "%s"
  namespace_name = "%s" 
  repository_name = "%s"
  image_version = "v2"
}

`
