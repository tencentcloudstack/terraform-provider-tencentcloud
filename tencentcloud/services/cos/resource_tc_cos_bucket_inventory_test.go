package cos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBucketInventoryResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBucketInventory,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_bucket_inventory.bucket_inventory", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cos_bucket_inventory.bucket_inventory",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCosBucketInventory = `
resource "tencentcloud_cos_bucket_inventory" "bucket_inventory" {
    name = "test123"
    bucket = "keep-test-1308919341"
    is_enabled = "true"
    included_object_versions = "Current"
    optional_fields {
        fields = ["Size", "ETag"]
    }
    filter {
        period {
            start_time = "1687276800"
        }
    }
    schedule {
        frequency = "Weekly"
    }
    destination {
        bucket = "qcs::cos:ap-guangzhou::keep-test-1308919341"
        account_id = ""
        format = "CSV"
        prefix = "cos_bucket_inventory"

    }
}
`
