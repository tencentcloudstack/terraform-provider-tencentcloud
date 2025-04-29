package vpn_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcvpc "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/vpc"

	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

// go test -i; go test -test.run TestAccTencentCloudVpnConnectionResource_basic -v
func TestAccTencentCloudVpnConnectionResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckVpnConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("tencentcloud_vpn_connection.connection"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "name", "vpn_connection_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "pre_share_key", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_encry_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_authen_algorithm", "MD5"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_local_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_remote_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_dh_group_name", "GROUP1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_exchange_mode", "MAIN"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_sa_lifetime_seconds", "86400"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_version", "IKEV1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_encrypt_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_integrity_algorithm", "MD5"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_seconds", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_pfs_dh_group", "DH-GROUP1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_traffic", "2560"),
					//resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "security_group_policy.0.remote_cidr_block.0", "3.3.3.0/32"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "net_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "encrypt_proto"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "route_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "vpn_proto"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_timeout", "30"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_action", "clear"),
				),
			},
			{
				Config: testAccVpnConnectionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("tencentcloud_vpn_connection.connection"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "name", "vpn_connection_test2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "pre_share_key", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "tags.test", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_encry_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_authen_algorithm", "SHA"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_local_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_remote_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_dh_group_name", "GROUP2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_exchange_mode", "AGGRESSIVE"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_sa_lifetime_seconds", "86401"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_version", "IKEV2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_encrypt_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_integrity_algorithm", "SHA1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_pfs_dh_group", "NULL"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_traffic", "2570"),
					//resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "security_group_policy.0.remote_cidr_block.0", "3.3.3.0/26"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "net_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "encrypt_proto"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "route_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "vpn_proto"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_enable", "0"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_timeout", "30"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_action", "clear"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "enable_health_check", "false"),
				),
			},
			{
				Config: testAccVpnConnectionConfigUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("tencentcloud_vpn_connection.connection"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "name", "vpn_connection_test2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "pre_share_key", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "tags.test", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_encry_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_authen_algorithm", "SHA"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_local_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_remote_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_dh_group_name", "GROUP2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_exchange_mode", "AGGRESSIVE"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_sa_lifetime_seconds", "86401"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_encrypt_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_integrity_algorithm", "SHA1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_pfs_dh_group", "NULL"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_traffic", "2570"),
					//resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "security_group_policy.0.remote_cidr_block.0", "3.3.3.0/26"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "net_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "encrypt_proto"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "route_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "vpn_proto"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_timeout", "40"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_action", "restart"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "enable_health_check", "false"),
				),
			},
			{
				Config: testAccVpnConnectionConfigUpdate3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("tencentcloud_vpn_connection.connection"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "name", "vpn_connection_test2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "pre_share_key", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "tags.test", "testt"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_encry_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_proto_authen_algorithm", "SHA"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_local_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_remote_identity", "ADDRESS"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_dh_group_name", "GROUP2"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_exchange_mode", "AGGRESSIVE"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ike_sa_lifetime_seconds", "86401"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_encrypt_algorithm", "3DES-CBC"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_integrity_algorithm", "SHA1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_pfs_dh_group", "NULL"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "ipsec_sa_lifetime_traffic", "2570"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "net_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "state"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "encrypt_proto"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "route_type", "StaticRoute"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpn_connection.connection", "vpn_proto"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_enable", "1"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_timeout", "40"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "dpd_action", "restart"),
					resource.TestCheckResourceAttr("tencentcloud_vpn_connection.connection", "enable_health_check", "false"),
				),
			},
		},
	})
}

func testAccCheckVpnConnectionDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpn_connection" {
			continue
		}
		request := vpc.NewDescribeVpnConnectionsRequest()
		request.VpnConnectionIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnConnectionsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnConnections(request)
			if e != nil {
				ee, ok := e.(*errors.TencentCloudSDKError)
				if !ok {
					return tccommon.RetryError(e)
				}
				if ee.Code == svcvpc.VPCNotFound {
					log.Printf("[CRITAL]%s api[%s] success, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return resource.NonRetryableError(e)
				} else {
					log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
						logId, request.GetAction(), request.ToJsonString(), e.Error())
					return tccommon.RetryError(e)
				}
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN connection failed, reason:%s\n", logId, err.Error())
			ee, ok := err.(*errors.TencentCloudSDKError)
			if !ok {
				return err
			}
			if ee.Code == "ResourceNotFound" {
				return nil
			} else {
				return err
			}
		} else {
			if len(response.Response.VpnConnectionSet) != 0 {
				return fmt.Errorf("VPN connection id is still exists")
			}
		}

	}
	return nil
}

func testAccCheckVpnConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("VPN connection instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("VPN connection id is not set")
		}
		conn := tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn()
		request := vpc.NewDescribeVpnConnectionsRequest()
		request.VpnConnectionIds = []*string{&rs.Primary.ID}
		var response *vpc.DescribeVpnConnectionsResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := conn.UseVpcClient().DescribeVpnConnections(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), e.Error())
				return tccommon.RetryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s read VPN connection failed, reason:%s\n", logId, err.Error())
			return err
		}
		if len(response.Response.VpnConnectionSet) != 1 {
			return fmt.Errorf("VPN connection id is not found")
		}
		return nil
	}
}

const testAccVpnConnectionConfig = `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.3.3.3"

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
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "MD5"
  ipsec_sa_lifetime_seconds  = 3600
  ipsec_pfs_dh_group         = "DH-GROUP1"
  ipsec_sa_lifetime_traffic  = 2560
  dpd_enable                 = 1
  dpd_timeout                = "30"
  dpd_action                 = "clear"
  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["3.3.3.0/32", ]
  }
  tags = {
    test = "test"
  }
}
`

const testAccVpnConnectionConfigUpdate = `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.3.3.3"
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
  name                       = "vpn_connection_test2"
  vpc_id                     = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.vpn.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.cgw.id
  pre_share_key              = "testt"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = tencentcloud_vpn_gateway.vpn.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.cgw.public_ip_address
  ike_dh_group_name          = "GROUP2"
  ike_version                = "IKEV2"
  ike_sa_lifetime_seconds    = 86401
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_sa_lifetime_seconds  = 7200
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_traffic  = 2570
  dpd_enable = 0
  dpd_timeout = "30"
  dpd_action = "clear"

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["3.3.3.0/26", ]
  }
  tags = {
    test = "testt"
  }
  enable_health_check = false
}
`

const testAccVpnConnectionConfigUpdate2 = `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.3.3.3"
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
  name                       = "vpn_connection_test2"
  vpc_id                     = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.vpn.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.cgw.id
  pre_share_key              = "testt"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = tencentcloud_vpn_gateway.vpn.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.cgw.public_ip_address
  ike_dh_group_name          = "GROUP2"
  ike_sa_lifetime_seconds    = 86401
  ike_version                = "IKEV2"
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_sa_lifetime_seconds  = 7200
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_traffic  = 2570
  dpd_enable = 1
  dpd_timeout = "40"
  dpd_action = "restart"

  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["3.3.3.0/26", ]
  }
  tags = {
    test = "testt"
  }
  enable_health_check = false
}
`

const testAccVpnConnectionConfigUpdate3 = `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test"
  public_ip_address = "1.3.3.3"
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
  name                       = "vpn_connection_test2"
  vpc_id                     = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  vpn_gateway_id             = tencentcloud_vpn_gateway.vpn.id
  customer_gateway_id        = tencentcloud_vpn_customer_gateway.cgw.id
  pre_share_key              = "testt"
  ike_proto_encry_algorithm  = "3DES-CBC"
  ike_proto_authen_algorithm = "SHA"
  ike_local_identity         = "ADDRESS"
  ike_exchange_mode          = "AGGRESSIVE"
  ike_local_address          = tencentcloud_vpn_gateway.vpn.public_ip_address
  ike_remote_identity        = "ADDRESS"
  ike_remote_address         = tencentcloud_vpn_customer_gateway.cgw.public_ip_address
  ike_dh_group_name          = "GROUP2"
  ike_sa_lifetime_seconds    = 86401
  ike_version                = "IKEV2"
  ipsec_encrypt_algorithm    = "3DES-CBC"
  ipsec_integrity_algorithm  = "SHA1"
  ipsec_sa_lifetime_seconds  = 7200
  ipsec_pfs_dh_group         = "NULL"
  ipsec_sa_lifetime_traffic  = 2570
  dpd_enable = 1
  dpd_timeout = "40"
  dpd_action = "restart"
  route_type = "StaticRoute"

  tags = {
    test = "testt"
  }
  enable_health_check = false
}
`
