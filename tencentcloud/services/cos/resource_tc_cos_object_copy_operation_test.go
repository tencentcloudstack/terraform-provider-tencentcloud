package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosObjectCopyOperationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosObjectCopyOperation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_object_copy_operation.object_copy", "id"),
				),
			},
		},
	})
}

const testAccCosObjectCopyOperation = `
resource "tencentcloud_cos_object_copy_operation" "object_copy" {
    bucket = "keep-copy-1308919341"
    key = "copy-acl.txt"
    source_url = "keep-test-1308919341.cos.ap-guangzhou.myqcloud.com/acl.txt"
}
`
