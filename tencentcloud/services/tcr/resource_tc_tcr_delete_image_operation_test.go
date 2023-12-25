package tcr_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixTcrDeleteImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrDeleteImageOperation, tcacctest.DefaultTCRInstanceId, tcacctest.DefaultTCRNamespace, tcacctest.DefaultTCRRepoName),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_delete_image_operation.delete_image_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "registry_id", tcacctest.DefaultTCRInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "namespace_name", tcacctest.DefaultTCRNamespace),
					resource.TestCheckResourceAttr("tencentcloud_tcr_delete_image_operation.delete_image_operation", "repository_name", tcacctest.DefaultTCRRepoName),
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
