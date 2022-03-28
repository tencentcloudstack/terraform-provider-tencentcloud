package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudVpnGatewayRoute_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnGatewayRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testVpnGatewayRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayRouteExists("tencentcloud_vpn_gateway_route.route1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "destination_cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "instance_type", "VPNCONN"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "priority", "100"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "status", "ENABLE"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_route.route1", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_route.route1", "route_id"),
				),
			},
			{
				Config: testVpnGatewayRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayRouteExists("tencentcloud_vpn_gateway_route.route1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "destination_cidr_block", "10.0.0.0/16"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "instance_type", "VPNCONN"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "priority", "100"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway_route.route1", "status", "DISABLE"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_route.route1", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway_route.route1", "route_id"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayRouteDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_gateway_route" {
			continue
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		err, result := vpcService.DescribeVpnGatewayRoutes(ctx, ids[0], nil)
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway route failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*errors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == VPCNotFound {
				return nil
			} else {
				return err
			}
		} else {
			if len(result) != 0 {
				return fmt.Errorf("VPN gateway route id is still exists")
			}
		}
	}
	return nil
}

func testAccCheckVpnGatewayRouteExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		vpcService := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN gateway route instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN gateway route id is not set")
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		err, result := vpcService.DescribeVpnGatewayRoutes(ctx, ids[0], nil)
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(result) != 1 {
			return fmt.Errorf("VPN gateway route id is not found")
		}
		return nil
	}
}

const testVpnGatewayRouteCreate = defaultVpnDataSource + `
# Create VPC
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway_route" "route1" {
  vpn_gateway_id = data.tencentcloud_vpn_gateways.foo.gateway_list.0.id
  destination_cidr_block = "10.0.0.0/16"
  instance_type = "VPNCONN"
  instance_id = data.tencentcloud_vpn_connections.conns.connection_list.0.id
  priority = "100"
  status = "ENABLE"
}
`
const testVpnGatewayRouteUpdate = defaultVpnDataSource + `
# Create VPC
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway_route" "route1" {
  vpn_gateway_id = data.tencentcloud_vpn_gateways.foo.gateway_list.0.id
  destination_cidr_block = "10.0.0.0/16"
  instance_type = "VPNCONN"
  instance_id = data.tencentcloud_vpn_connections.conns.connection_list.0.id
  priority = "100"
  status = "DISABLE"
}
`
