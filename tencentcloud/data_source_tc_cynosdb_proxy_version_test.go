package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbProxyVersionDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_proxy_version.proxy_version"),
				),
			},
		},
	})
}

const testAccCynosdbProxyVersionDataSource = `
data "tencentcloud_cynosdb_proxy_version" "proxy_version" {
  cluster_id     = "cynosdbmysql-bws8h88b"
}
`
