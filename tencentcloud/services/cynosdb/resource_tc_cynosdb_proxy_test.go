package cynosdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbProxyResource_basic -v
func TestAccTencentCloudCynosdbProxyResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckCynosdbProxyDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbProxyExists("tencentcloud_cynosdb_proxy.proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "cpu"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "mem"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "unique_vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "unique_subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "connection_pool_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "open_connection_pool"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "connection_pool_time_out"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "security_group_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "description"),
				),
			},
			{
				Config: testAccCynosdbProxyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbProxyExists("tencentcloud_cynosdb_proxy.proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "cpu"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "mem"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "unique_vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "unique_subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "connection_pool_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "open_connection_pool"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "connection_pool_time_out"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "security_group_ids.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy.proxy", "description"),
				),
			},
		},
	})
}

func testAccCheckCynosdbProxyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb proxy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb proxy id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]

		proxy, err := cynosdbService.DescribeCynosdbProxyById(ctx, clusterId, "")
		if err != nil {
			return err
		}

		if proxy.ProxyGroupInfos == nil || *proxy.TotalCount == 0 {
			return fmt.Errorf("cynosdb proxy doesn't exists: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckCynosdbProxyDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_proxy" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]

		proxy, err := cynosdbService.DescribeCynosdbProxyById(ctx, clusterId, "")
		if err != nil {
			return err
		}

		if proxy.ProxyGroupInfos == nil || *proxy.TotalCount == 0 {
			return nil
		}

		return fmt.Errorf("cynosdb proxy still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccCynosdbProxy = `
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = "vpc-k1t8ickr"
  unique_subnet_id         = "subnet-jdi5xn22"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc sample"
  proxy_zones {
    proxy_node_zone  = "ap-guangzhou-7"
    proxy_node_count = 2
  }
}
`

const testAccCynosdbProxyUpdate = `
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  cpu                      = 4
  mem                      = 6000
  unique_vpc_id            = "vpc-4owdpnwr" 
  unique_subnet_id         = "subnet-dwj7ipnc"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc sample test"
  proxy_zones {
    proxy_node_zone  = "ap-guangzhou-7"
    proxy_node_count = 2
  }
}
`
