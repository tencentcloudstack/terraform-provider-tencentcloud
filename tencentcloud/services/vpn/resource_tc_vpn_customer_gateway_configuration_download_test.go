package vpn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudVpnCustomerGatewayConfigurationDownloadResource_basic -v -timeout=0
func TestAccTencentCloudVpnCustomerGatewayConfigurationDownloadResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcVpnCustomerGatewayConfigurationDownload,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpn_customer_gateway_configuration_download.vpn_customer_gateway_configuration_download", "id")),
			},
		},
	})
}

const testAccVpcVpnCustomerGatewayConfigurationDownload = `
resource "tencentcloud_vpn_customer_gateway" "cgw" {
  name              = "terraform_test_download"
  public_ip_address = "10.3.1.18"

}

# Create VPC and Subnet
data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "vpn" {
  name      = "terraform_test_download"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}
resource "tencentcloud_vpn_connection" "connection" {
  name                       = "test_download"
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
  dpd_enable = 1
  dpd_timeout = "30"
  dpd_action = "clear"
  security_group_policy {
    local_cidr_block  = "172.16.0.0/16"
    remote_cidr_block = ["3.3.3.0/32", ]
  }
  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id    = tencentcloud_vpn_gateway.vpn.id
  vpn_connection_id = tencentcloud_vpn_connection.connection.id
  customer_gateway_vendor {
    platform         = "comware"
    software_version = "V1.0"
    vendor_name      = "h3c"
  }
  interface_name    = "test"
}

`
