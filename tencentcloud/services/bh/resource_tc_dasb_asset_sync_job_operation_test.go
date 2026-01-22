package bh_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixDasbAssetSyncJobOperationOperationResource_basic -v
func TestAccTencentCloudNeedFixDasbAssetSyncJobOperationOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDasbAssetSyncJobOperationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_asset_sync_job_operation.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dasb_asset_sync_job_operation.example", "category"),
				),
			},
		},
	})
}

const testAccDasbAssetSyncJobOperationOperation = `
resource "tencentcloud_dasb_asset_sync_job_operation" "example" {
  category = 1
}
`
