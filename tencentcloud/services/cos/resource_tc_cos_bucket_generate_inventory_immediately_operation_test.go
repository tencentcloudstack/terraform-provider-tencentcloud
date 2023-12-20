package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBucketGenerateInventoryImmediatelyOperationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketGenerateInventoryImmediatelyOperation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_generate_inventory_immediately_operation.generate_inventory_immediately", "id"),
				),
			},
		},
	})
}

const testAccCosBucketGenerateInventoryImmediatelyOperation = `
resource "tencentcloud_cos_bucket_generate_inventory_immediately_operation" "generate_inventory_immediately" {
    inventory_id = "test"
    bucket = "keep-test-1308919341"
}
`
