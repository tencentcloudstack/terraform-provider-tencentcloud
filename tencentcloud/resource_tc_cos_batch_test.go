package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCosBatchResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCosBatchResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cos_batch.cos_batch", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cos_batch.cos_batch",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCosBatchResource = `
resource "tencentcloud_cos_batch" "cos_batch" {
    uin = "100022975249"
    appid = "1308919341"
    confirmation_required = true
    description = "cos_batch"
    priority = 1
	status = "Cancelled"
    role_arn = "qcs::cam::uin/100022975249:roleName/COSBatch_QCSRole"
    manifest {
        location {
            etag = "64357de8fd75a3abae2200135a2c9627"
            object_arn = "qcs::cos:ap-guangzhou:uid/1308919341:keep-test-1308919341/cos_bucket_inventory/1308919341/keep-test/test/20230621/manifest.json"
        }
        spec {
            format = "COSInventoryReport_CSV_V1"
        }
    }
    operation {
        cos_put_object_copy {
            access_control_directive = "Copy"
            metadata_directive = "Copy"
            prefix_replace = false
            storage_class = "STANDARD"
            tagging_directive = "Copy"
            target_resource = "qcs::cos:ap-guangzhou:uid/1308919341:cos-lock-1308919341"
        }
    }
    report {
        bucket = "qcs::cos:ap-guangzhou:uid/1308919341:keep-test-1308919341"
        enabled = "true"
        format = "Report_CSV_V1"
        report_scope = "AllTasks"
    }
}
`
