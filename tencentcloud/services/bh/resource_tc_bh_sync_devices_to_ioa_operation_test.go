package bh_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBhSyncDevicesToIoaOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccBhSyncDevicesToIoaOperation,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_bh_sync_devices_to_ioa_operation.example", "id"),
			),
		}},
	})
}

const testAccBhSyncDevicesToIoaOperation = `
resource "tencentcloud_bh_sync_devices_to_ioa_operation" "example" {
  device_id_set = [
    1934,
    1964,
    1895,
  ]
}
`
