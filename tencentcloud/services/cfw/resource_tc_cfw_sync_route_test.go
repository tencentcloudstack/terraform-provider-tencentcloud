package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwSyncRouteResource_basic -v
func TestAccTencentCloudNeedFixCfwSyncRouteResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwSyncRoute,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_sync_route.example", "id"),
				),
			},
		},
	})
}

const testAccCfwSyncRoute = `
resource "tencentcloud_cfw_sync_route" "example" {
  sync_type = "Route"
  fw_type   = "nat"
}
`
