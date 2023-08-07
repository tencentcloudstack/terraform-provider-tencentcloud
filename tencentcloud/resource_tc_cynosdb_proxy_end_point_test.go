package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbProxyEndPointResource_basic -v
func TestAccTencentCloudCynosdbProxyEndPointResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckCynosdbProxyEndPointDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbProxyEndPoint,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbProxyEndPointExists("tencentcloud_cynosdb_proxy_end_point.proxy_end_point"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy_end_point.proxy_end_point", "id"),
				),
			},
			{
				Config: testAccCynosdbProxyEndPointUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbProxyEndPointExists("tencentcloud_cynosdb_proxy_end_point.proxy_end_point"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_proxy_end_point.proxy_end_point", "id"),
				),
			},
		},
	})
}

func testAccCheckCynosdbProxyEndPointDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := CynosdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_proxy_end_point" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		proxyGroupId := idSplit[1]

		proxyEndPoint, err := service.DescribeCynosdbProxyEndPointById(ctx, clusterId, proxyGroupId)
		if err != nil {
			return err
		}

		if proxyEndPoint == nil {
			return nil
		}

		return fmt.Errorf("cynosdb proxy end point still exists: %s", rs.Primary.ID)
	}

	return nil
}

func testAccCheckCynosdbProxyEndPointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb proxy end point %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb proxy end point id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		clusterId := idSplit[0]
		proxyGroupId := idSplit[1]

		service := CynosdbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		proxyEndPoint, err := service.DescribeCynosdbProxyEndPointById(ctx, clusterId, proxyGroupId)
		if err != nil {
			return err
		}

		if proxyEndPoint == nil {
			return fmt.Errorf("cynosdb proxy end point doesn't exist: %s", rs.Primary.ID)
		}

		return nil
	}
}

const testAccCynosdbProxyEndPoint = CommonCynosdb + DefaultCrsVar + `
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = var.cynosdb_cluster_id
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = var.vpc_id
  unique_subnet_id         = var.subnet_id
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

resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = tencentcloud_cynosdb_proxy.proxy.cluster_id
  unique_vpc_id            = var.vpc_id
  unique_subnet_id         = var.subnet_id
  vip                      = "172.16.96.128"
  vport                    = "3306"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc value"
  weight_mode              = "system"
  auto_add_ro              = "yes"
  fail_over                = "yes"
  consistency_type         = "global"
  rw_type                  = "READWRITE"
  consistency_time_out     = 30
  trans_split              = true
  access_mode              = "nearby"
  instance_weights {
    instance_id = tencentcloud_cynosdb_proxy.proxy.ro_instances.0.instance_id
    weight      = 1
  }
}
`

const testAccCynosdbProxyEndPointUpdate = CommonCynosdb + DefaultCrsVar + `
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = var.cynosdb_cluster_id
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = var.vpc_id
  unique_subnet_id         = var.subnet_id
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

resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id               = tencentcloud_cynosdb_proxy.proxy.cluster_id
  unique_vpc_id            = var.vpc_id
  unique_subnet_id         = var.subnet_id
  vip                      = "172.16.96.158"
  vport                    = "3306"
  open_connection_pool     = "no"
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc value"
  weight_mode              = "system"
  auto_add_ro              = "no"
  rw_type                  = "READWRITE"
  trans_split              = true
  access_mode              = "balance"
  instance_weights {
    instance_id = tencentcloud_cynosdb_proxy.proxy.ro_instances.0.instance_id
    weight      = 1
  }
}
`
