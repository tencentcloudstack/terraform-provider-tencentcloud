package gwlb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGwlbInstanceAssociateTargetGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGwlbInstanceAssociateTargetGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_gwlb_instance_associate_target_group.gwlb_instance_associate_target_groups", "id")),
			},
		},
	})
}

const testAccGwlbInstanceAssociateTargetGroup = `

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_gwlb_instance" "gwlb_instance" {
  vpc_id = tencentcloud_vpc.vpc.id
  subnet_id = tencentcloud_subnet.subnet.id
  load_balancer_name = "tf-test"
  lb_charge_type = "POSTPAID_BY_HOUR"
  tags {
    tag_key = "test_key"
    tag_value = "tag_value"
  }
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

resource "tencentcloud_gwlb_instance_associate_target_group" "gwlb_instance_associate_target_groups" {
  load_balancer_id = tencentcloud_gwlb_instance.gwlb_instance.id
  target_group_id = tencentcloud_gwlb_target_group.gwlb_target_group.id
}
`
