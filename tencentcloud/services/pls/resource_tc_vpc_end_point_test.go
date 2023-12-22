package pls_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudVpcEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEndPoint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_end_point.end_point",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcEndPoint = `

resource "tencentcloud_vpc_end_point" "end_point" {
  vpc_id = "vpc-391sv4w3"
  subnet_id = "subnet-ljyn7h30"
  end_point_name = "terraform-test"
  end_point_service_id = "vpcsvc-98jddhcz"
  end_point_vip = "10.0.2.2"
}

`
