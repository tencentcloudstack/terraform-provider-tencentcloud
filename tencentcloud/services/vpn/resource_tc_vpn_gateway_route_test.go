package vpn_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpn"

	"github.com/stretchr/testify/assert"
)

// go test -i; go test -test.run TestAccTencentCloudVpnGatewayRoute_basic -v -timeout=0
func TestAccTencentCloudVpnGatewayRoute_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_gateway_route" {
			continue
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		err, result := vpcService.DescribeVpnGatewayRoutes(ctx, ids[0], nil)
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway route failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == svcvpc.VPCNotFound {
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		vpcService := svcvpc.NewVpcService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN gateway route instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN gateway route id is not set")
		}
		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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

const testVpnGatewayRouteCreate = tcacctest.DefaultVpnDataSource + `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.14.14.14"

}

# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "vpn" {
  name      = "terraform_update"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_connection" "connection" {
  name                       = "vpn_connection_test"
  vpc_id                     = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.vpn.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.cgw.id
  pre_share_key              = "test"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "MD5"
  ike_local_identity         = "ADDRESS"
  ike_local_address          = tencentcloud_vpn_gateway.vpn.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.cgw.public_ip_address
  ike_dh_group_name          = "GROUP1"
  ike_sa_lifetime_seconds    = 86400
  ike_version                = "IKEV1"
  route_type                 = "StaticRoute"
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "MD5"
  ipsec_sa_lifetime_seconds  = 3600
  ipsec_pfs_dh_group         = "DH-GROUP1"
  ipsec_sa_lifetime_traffic  = 2560
  dpd_enable                 = 1
  dpd_timeout                = "30"
  dpd_action                 = "clear"
  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_gateway_route" "route1" {
  vpn_gateway_id         = tencentcloud_vpn_gateway.vpn.id
  destination_cidr_block = "10.0.0.0/16"
  instance_type          = "VPNCONN"
  instance_id            = tencentcloud_vpn_connection.connection.id
  priority               = "100"
  status                 = "ENABLE"
}
`
const testVpnGatewayRouteUpdate = tcacctest.DefaultVpnDataSource + `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.14.14.14"

}

# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "vpn" {
  name      = "terraform_update"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_connection" "connection" {
  name                       = "vpn_connection_test"
  vpc_id                     = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.vpn.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.cgw.id
  pre_share_key              = "test"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "MD5"
  ike_local_identity         = "ADDRESS"
  ike_local_address          = tencentcloud_vpn_gateway.vpn.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.cgw.public_ip_address
  ike_dh_group_name          = "GROUP1"
  ike_sa_lifetime_seconds    = 86400
  ike_version                = "IKEV1"
  route_type                 = "StaticRoute"
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "MD5"
  ipsec_sa_lifetime_seconds  = 3600
  ipsec_pfs_dh_group         = "DH-GROUP1"
  ipsec_sa_lifetime_traffic  = 2560
  dpd_enable                 = 1
  dpd_timeout                = "30"
  dpd_action                 = "clear"
  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_gateway_route" "route1" {
  vpn_gateway_id         = tencentcloud_vpn_gateway.vpn.id
  destination_cidr_block = "10.0.0.0/16"
  instance_type          = "VPNCONN"
  instance_id            = tencentcloud_vpn_connection.connection.id
  priority               = "100"
  status                 = "DISABLE"
}
`

type mockMetaForVpnGatewayRoute struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForVpnGatewayRoute) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForVpnGatewayRoute{}

func newMockMetaForVpnGatewayRoute() *mockMetaForVpnGatewayRoute {
	return &mockMetaForVpnGatewayRoute{client: &connectivity.TencentCloudClient{}}
}

// TestVpnGatewayRouteServiceDescribe_ResourceNotFound verifies that when the
// cloud API DescribeVpnGatewayRoutes returns a ResourceNotFound error, the
// service method returns (nil, nil) instead of propagating the error, so the
// upper Read path can treat the route as deleted and clear the id.
func TestVpnGatewayRouteServiceDescribe_ResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaForVpnGatewayRoute().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeVpnGatewayRoutes", func(request *vpc.DescribeVpnGatewayRoutesRequest) (*vpc.DescribeVpnGatewayRoutesResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "ResourceNotFound",
			Message: "the resource not found",
		}
	})

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcvpc.NewVpcService(newMockMetaForVpnGatewayRoute().GetAPIV3Conn())

	err, result := service.DescribeVpnGatewayRoutes(ctx, "vpngw-test", nil)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// TestVpnGatewayRouteRead_ResourceNotFound verifies that the resource Read
// clears d.Id() and returns nil when the cloud API returns ResourceNotFound,
// which is the behavior expected after the service-layer fix.
func TestVpnGatewayRouteRead_ResourceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	vpcClient := &vpc.Client{}
	patches.ApplyMethodReturn(newMockMetaForVpnGatewayRoute().client, "UseVpcClient", vpcClient)

	patches.ApplyMethodFunc(vpcClient, "DescribeVpnGatewayRoutes", func(request *vpc.DescribeVpnGatewayRoutesRequest) (*vpc.DescribeVpnGatewayRoutesResponse, error) {
		return nil, &sdkErrors.TencentCloudSDKError{
			Code:    "ResourceNotFound",
			Message: "the resource not found",
		}
	})

	meta := newMockMetaForVpnGatewayRoute()
	res := vpn.ResourceTencentCloudVpnGatewayRoute()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"vpn_gateway_id":         "vpngw-test",
		"destination_cidr_block": "10.0.0.0/16",
		"instance_type":          "VPNCONN",
		"instance_id":            "vpnx-test",
		"priority":               100,
		"status":                 "ENABLE",
	})
	d.SetId(helper.IdFormat("vpngw-test", "route-test"))

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}
