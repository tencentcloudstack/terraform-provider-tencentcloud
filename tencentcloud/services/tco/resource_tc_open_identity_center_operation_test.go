package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOpenIdentityCenterOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenIdentityCenterOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_open_identity_center_operation.open_identity_center_operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_open_identity_center_operation.open_identity_center_operation", "zone_id"),
				),
			},
		},
	})
}

const testAccOpenIdentityCenterOperation = `
resource "tencentcloud_open_identity_center_operation" "open_identity_center_operation" {
    zone_name = "test"
}
`
