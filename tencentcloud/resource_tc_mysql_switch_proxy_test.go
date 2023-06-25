package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlSwitchProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSwitchProxy,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_switch_proxy.switch_proxy", "id")),
			},
		},
	})
}

const testAccMysqlSwitchProxy = `

resource "tencentcloud_mysql_switch_proxy" "switch_proxy" {
  instance_id = ""
  proxy_group_id = ""
}

`
