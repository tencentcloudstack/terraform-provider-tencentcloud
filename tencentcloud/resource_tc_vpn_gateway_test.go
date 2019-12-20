package tencentcloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	errors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func TestAccTencentCloudVpnGateway_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("tencentcloud_vpn_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "tags.test", "tf"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway.my_cgw", "state"),
				),
			},
			{
				Config: testAccVpnGatewayConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("tencentcloud_vpn_gateway.my_cgw"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "name", "terraform_update"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "bandwidth", "5"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_gateway.my_cgw", "tags.test", "test"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_gateway.my_cgw", "state"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)

	conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_gateway" {
			continue
		}
		request := vpc.NewDescribeVpnGatewaysRequest()
		request.VpnGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnGateways(request)
			if e != nil {
				ee, ok := e.(*errors.TencentCloudSDKError)
				if !ok {
					return retryError(e)
				}
				if ee.Code == "ResourceNotFound" {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return retryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
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
			if len(response.Response.VpnGatewaySet) != 0 {
				return fmt.Errorf("VPN gateway id is still exists")
			}
		}

	}
	return nil
}

func testAccCheckVpnGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN gateway instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN gateway id is not set")
		}
		conn := testAccProvider.Meta().(*TencentCloudClient).apiV3Conn
		request := vpc.NewDescribeVpnGatewaysRequest()
		request.VpnGatewayIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnGatewaysResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnGateways(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN gateway failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.VpnGatewaySet) != 1 {
			return fmt.Errorf("VPN gateway id is not found")
		}
		return nil
	}
}

const testAccVpnGatewayConfig = `
# Create VPC
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_test"
  vpc_id    = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
  bandwidth = 10
  zone      = "ap-guangzhou-3"

  tags = {
    test = "tf"
  }
}
`
const testAccVpnGatewayConfigUpdate = `
# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_update"
  vpc_id    = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}

`
