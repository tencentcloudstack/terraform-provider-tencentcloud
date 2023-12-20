package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbProxyNodeDataSource_basic -v
func TestAccTencentCloudCynosdbProxyNodeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckCynosdbProxyDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyNodeDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_proxy_node.proxy_node"),
				),
			},
		},
	})
}

const testAccCynosdbProxyNodeDataSource = testAccCynosdbProxy + `
data "tencentcloud_cynosdb_proxy_node" "proxy_node" {
  order_by      = "CREATETIME"
  order_by_type = "DESC"
  filters {
    names       = ["ClusterId"]
    values      = [tencentcloud_cynosdb_proxy.proxy.id]
    exact_match = false
    name        = "ClusterId"
  }
}
`
