package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixPrivateDnsExtendEndPointResource_basic -v
func TestAccTencentCloudNeedFixPrivateDnsExtendEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsExtendEndPoint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_extend_end_point.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_extend_end_point.example", "end_point_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_extend_end_point.example", "end_point_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_extend_end_point.example", "forward_ip.#"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_extend_end_point.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsExtendEndPoint = `
resource "tencentcloud_private_dns_extend_end_point" "example" {
  end_point_name   = "tf-example"
  end_point_region = "ap-jakarta"
  forward_ip {
    access_type       = "CCN"
    host              = "1.1.1.1"
    port              = 8080
    vpc_id            = "vpc-2qjckjg2"
    access_gateway_id = "ccn-eo13f8ub"
  }
}
`
