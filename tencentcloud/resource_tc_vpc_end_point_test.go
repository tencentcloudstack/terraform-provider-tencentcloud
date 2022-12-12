package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudVpcEndPointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
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
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip = "10.0.2.1"
}

`
