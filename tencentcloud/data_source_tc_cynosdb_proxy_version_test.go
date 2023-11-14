package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbProxyVersionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyVersionDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_proxy_version.proxy_version")),
			},
		},
	})
}

const testAccCynosdbProxyVersionDataSource = `

data "tencentcloud_cynosdb_proxy_version" "proxy_version" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  proxy_group_id = "无"
    }

`
