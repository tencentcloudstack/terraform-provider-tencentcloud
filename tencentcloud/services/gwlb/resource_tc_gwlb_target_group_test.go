package gwlb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGwlbTargetGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGwlbTargetGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_gwlb_target_group.gwlb_target_group", "id")),
			},
			{
				Config: testAccGwlbTargetGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gwlb_target_group.gwlb_target_group", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gwlb_target_group.gwlb_target_group", "health_check.0.timeout", "5"),
				),
			},
			{
				ResourceName:      "tencentcloud_gwlb_target_group.gwlb_target_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGwlbTargetGroup = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_gwlb_target_group" "gwlb_target_group" {
  target_group_name = "tf-test"
  vpc_id = tencentcloud_vpc.vpc.id
  port = 6081
  health_check {
    health_switch = true
    protocol = "tcp"
    port = 6081
    timeout = 2
    interval_time = 5
    health_num = 3
    un_health_num = 3
  }
}
`

const testAccGwlbTargetGroupUpdate = `
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_gwlb_target_group" "gwlb_target_group" {
  target_group_name = "tf-test"
  vpc_id = tencentcloud_vpc.vpc.id
  port = 6081
  health_check {
    health_switch = true
    protocol = "tcp"
    port = 6081
    timeout = 5
    interval_time = 5
    health_num = 3
    un_health_num = 3
  }
}
`
