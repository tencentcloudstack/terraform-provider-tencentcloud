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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_end_point.end_point", "id"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "security_groups_ids.0", "sg-ghvp9djf"),
				),
			},

			{
				Config: testAccVpcEndPointUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "end_point_name", "terraform_test_for"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_end_point.end_point", "security_groups_ids.0", "sg-3k7vtgf7"),
				),
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
  end_point_name       = "terraform_test"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  security_groups_ids  = [
    "sg-ghvp9djf",
    "sg-if748odn",
    "sg-3k7vtgf7",
  ]
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}

`

const testAccVpcEndPointUpdate = `

resource "tencentcloud_vpc_end_point" "end_point" {
  end_point_name       = "terraform_test_for"
  end_point_service_id = "vpcsvc-5y4yb85d"
  end_point_vip        = "10.0.0.58"
  security_groups_ids  = [
	"sg-3k7vtgf7",
    "sg-ghvp9djf",
    "sg-if748odn",
  ]
  subnet_id = "subnet-cpknsqgo"
  vpc_id    = "vpc-gmq0mxoj"
}

`
