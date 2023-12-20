package dc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixDcPublicDirectConnectTunnelRoutesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcPublicDirectConnectTunnelRoutesDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dc_public_direct_connect_tunnel_routes.public_direct_connect_tunnel_routes")),
			},
		},
	})
}

const testAccDcPublicDirectConnectTunnelRoutesDataSource = `

data "tencentcloud_dc_public_direct_connect_tunnel_routes" "public_direct_connect_tunnel_routes" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
}
`
