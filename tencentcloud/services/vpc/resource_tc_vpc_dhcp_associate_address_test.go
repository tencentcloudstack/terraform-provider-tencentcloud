package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcDhcpAssociateAddressResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDhcpAssociateAddress,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_dhcp_associate_address.dhcp_associate_address", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_dhcp_associate_address.dhcp_associate_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcDhcpAssociateAddress = `

resource "tencentcloud_vpc_dhcp_associate_address" "dhcp_associate_address" {
  dhcp_ip_id = ""
  address_ip = ""
}

`
