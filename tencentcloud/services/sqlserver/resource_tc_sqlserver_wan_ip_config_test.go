package sqlserver_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudSqlserverWanIpConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccSqlserverWanIpConfig,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_wan_ip_config.sqlserver_wan_ip_config", "id")),
		}, {
			ResourceName:      "tencentcloud_sqlserver_wan_ip_config.sqlserver_wan_ip_config",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccSqlserverWanIpConfig = `

resource "tencentcloud_sqlserver_wan_ip_config" "sqlserver_wan_ip_config" {
}
`
