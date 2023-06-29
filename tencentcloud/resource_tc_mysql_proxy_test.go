package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlProxyResource_basic -v
func TestAccTencentCloudMysqlProxyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlProxyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlProxyExists("tencentcloud_mysql_proxy.proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "uniq_vpc_id", defaultDcdbInsVpcId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "uniq_subnet_id", defaultDcdbInsIdSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.node_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.cpu", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.mem", "4000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "security_group.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "desc", "desc"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "connection_pool_limit", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "proxy_version"),
				),
			},
			{
				ResourceName:      "tencentcloud_mysql_proxy.proxy",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMysqlProxyUp,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMysqlProxyExists("tencentcloud_mysql_proxy.proxy"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_proxy.proxy", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "uniq_vpc_id", defaultDcdbInsVpcId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "uniq_subnet_id", defaultDcdbInsIdSubnetId),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.node_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.cpu", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.mem", "4000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "proxy_node_custom.0.zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "security_group.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "desc", "terraform test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "connection_pool_limit", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "vip", "172.16.151.23"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_proxy.proxy", "vport", "3306"),
				),
			},
		},
	})
}

func testAccCheckMysqlProxyDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_proxy" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		proxyGroupId := idSplit[1]
		// proxyAddressId := idSplit[2]

		instance, err := mysqlService.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
		if instance != nil {
			return fmt.Errorf("mysql proxy still exist")
		}
		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
	}
	return nil
}

func testAccCheckMysqlProxyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mysql proxy %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mysql proxy id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		proxyGroupId := idSplit[1]
		// proxyAddressId := idSplit[2]

		mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		instance, err := mysqlService.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
		if instance == nil {
			return fmt.Errorf("mysql proxy %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

const testAccMysqlProxyVar = `
variable "instance_id" {
  default = "` + defaultDbBrainInstanceId + `"
}

variable "vpc_id" {
	default = "` + defaultDcdbInsVpcId + `"
}

variable "subnet_id" {
	default = "` + defaultDcdbInsIdSubnetId + `"
}

`

const testAccMysqlProxy = testAccMysqlProxyVar + `

resource "tencentcloud_mysql_proxy" "proxy" {
	instance_id    = var.instance_id
	uniq_vpc_id    = var.vpc_id
	uniq_subnet_id = var.subnet_id
	proxy_node_custom {
	  node_count = 1
	  cpu        = 2
	  mem        = 4000
	  region     = "ap-guangzhou"
	  zone       = "ap-guangzhou-3"
	}
	security_group        = ["sg-edmur627"]
	desc                  = "desc"
	connection_pool_limit = 1
}

`

const testAccMysqlProxyUp = testAccMysqlProxyVar + `

resource "tencentcloud_mysql_proxy" "proxy" {
	instance_id    = var.instance_id
	uniq_vpc_id    = var.vpc_id
	uniq_subnet_id = var.subnet_id
	proxy_node_custom {
	  node_count = 1
	  cpu        = 2
	  mem        = 4000
	  region     = "ap-guangzhou"
	  zone       = "ap-guangzhou-3"
	}
	security_group        = ["sg-edmur627"]
	desc                  = "terraform test"
	connection_pool_limit = 2
	vip                   = "172.16.151.23"
	vport                 = 3306
}

`
