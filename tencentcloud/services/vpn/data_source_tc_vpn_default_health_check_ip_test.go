package vpn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpnDefaultHealthCheckIpDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpnDefaultHealthCheckIpDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpn_default_health_check_ip.default_health_check_ip")),
			},
		},
	})
}

const testAccVpnDefaultHealthCheckIpDataSource = `

data "tencentcloud_vpc_instances" "foo" {
  name = "Default-VPC"
}

resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "terraform_test_health_check"
  vpc_id    = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
  bandwidth = 10

  tags = {
    test = "tf"
  }
}

data "tencentcloud_vpn_default_health_check_ip" "default_health_check_ip" {
  vpn_gateway_id = tencentcloud_vpn_gateway.my_cgw.id
}

`
