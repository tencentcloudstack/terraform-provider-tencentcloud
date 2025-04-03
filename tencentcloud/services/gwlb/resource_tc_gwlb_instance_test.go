package gwlb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudGwlbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGwlbInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_gwlb_instance.gwlb_instance", "id")),
			},
			{
				Config: testAccGwlbInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gwlb_instance.gwlb_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gwlb_instance.gwlb_instance", "load_balancer_name", "tf-test-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_gwlb_instance.gwlb_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGwlbInstance = `
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
`

const testAccGwlbInstanceUpdate = `
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
  load_balancer_name = "tf-test-update"
  lb_charge_type = "POSTPAID_BY_HOUR"
  tags {
    tag_key = "test_key"
    tag_value = "tag_value"
  }
}
`
