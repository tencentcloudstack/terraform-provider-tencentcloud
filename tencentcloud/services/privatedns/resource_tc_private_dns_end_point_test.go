package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixPrivateDnsEndPointResource_basic -v
func TestAccTencentCloudNeedFixPrivateDnsEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateDnsEndPoint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_end_point.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_end_point.example", "end_point_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_end_point.example", "end_point_service_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_end_point.example", "end_point_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_private_dns_end_point.example", "ip_num"),
				),
			},
			{
				ResourceName:      "tencentcloud_private_dns_end_point.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivateDnsEndPoint = `
resource "tencentcloud_private_dns_end_point" "example" {
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-61wcwmar"
  end_point_region     = "ap-guangzhou"
  ip_num               = 1
}
`
