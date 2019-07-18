package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudClbInstancesDataSource_internal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstancesDataSource_internal,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_name", "tf-clb-internal"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.network_type", "INTERNAL"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_vips.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status"),
				),
			},
		},
	})
}

func TestAccTencentCloudClbInstancesDataSource_open(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClbInstancesDataSource_open,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckClbInstanceExists("tencentcloud_clb_instance.clb"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.network_type", "OPEN"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.clb_vips.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.vpc_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.status"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.target_region_info_region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_clb_instances.clbs", "clb_list.0.target_region_info_vpc"),
					resource.TestCheckResourceAttr("data.tencentcloud_clb_instances.clbs", "clb_list.0.security_groups.#", "1"),
				),
			},
		},
	})
}

const testAccClbInstancesDataSource_internal = `
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
	availability_zone="${var.availability_zone}"
	name="guagua-ci-temp-test"
	vpc_id="${tencentcloud_vpc.foo.id}"
	cidr_block="10.0.20.0/28"
	is_multicast=false
}

resource "tencentcloud_clb_instance" "clb" {
	network_type      = "INTERNAL"
	clb_name      = "tf-clb-internal"
  	vpc_id    		  = "${tencentcloud_vpc.foo.id}"
  	subnet_id		  = "${tencentcloud_subnet.subnet.id}"
	project_id = 0
	
}

data "tencentcloud_clb_instances" "clbs" {
	clb_id = "${tencentcloud_clb_instance.clb.id}"
}
`
const testAccClbInstancesDataSource_open = `
resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}
variable "availability_zone" {
	default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
    name="guagua-ci-temp-test"
    cidr_block="10.0.0.0/16"
}

resource "tencentcloud_clb_instance" "clb" {
	network_type      = "OPEN"
	clb_name      = "tf-clb-open"
	project_id = 0
	vpc_id="${tencentcloud_vpc.foo.id}"
	target_region_info_region = "ap-guangzhou"
	target_region_info_vpc_id = "${tencentcloud_vpc.foo.id}"
	}

	security_groups = ["${tencentcloud_security_group.foo.id}"]
}
data "tencentcloud_clb_instances" "clbs" {
	clb_id = "${tencentcloud_clb_instance.clb.id}"
}
`
