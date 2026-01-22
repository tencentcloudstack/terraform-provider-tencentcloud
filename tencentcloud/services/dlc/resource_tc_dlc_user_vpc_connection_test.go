package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcUserVpcConnectionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUserVpcConnection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_vpc_connection.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_vpc_connection.example", "user_vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_vpc_connection.example", "user_subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_vpc_connection.example", "user_vpc_endpoint_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_vpc_connection.example", "engine_network_id"),
				),
			},
		},
	})
}

const testAccDlcUserVpcConnection = `
resource "tencentcloud_dlc_user_vpc_connection" "example" {
  user_vpc_id            = "vpc-f7fa1fu5"
  user_subnet_id         = "subnet-ds2t3udw"
  user_vpc_endpoint_name = "tf-example"
  engine_network_id      = "DataEngine-Network-2mfg9icb"
}
`
