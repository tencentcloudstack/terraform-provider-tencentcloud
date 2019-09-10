package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDnatsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnatsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dnats.multi_dnats"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnats.multi_dnats", "dnat_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dnats.multi_dnats", "dnat_list.0.description", "test"),
				),
			},
		},
	})
}

const testAccTencentCloudDnatsDataSourceConfig_basic = `
data "tencentcloud_availability_zones" "my_favorate_zones" {
	name = "ap-guangzhou-3"
  }
  
  data "tencentcloud_image" "my_favorate_image" {
	filter {
	  name   = "image-type"
	  values = ["PUBLIC_IMAGE"]
	}
  }
  
  # Create VPC and Subnet
  data "tencentcloud_vpc_instances" "foo" {
	name = "Default-VPC"
}

data "tencentcloud_vpc_subnets" "foo" {
	subnet_id = "subnet-pqfek0t8"
  }
  
  # Create EIP 
  resource "tencentcloud_eip" "eip_dev_dnat" {
	name = "terraform_test"
  }
  resource "tencentcloud_eip" "eip_test_dnat" {
	name = "terraform_test"
  }
  
  # Create NAT Gateway
  resource "tencentcloud_nat_gateway" "my_nat" {
	vpc_id           = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
	name             = "terraform test"
	max_concurrent   = 3000000
	bandwidth        = 500
	assigned_eip_set = [
	  "${tencentcloud_eip.eip_dev_dnat.public_ip}",
	  "${tencentcloud_eip.eip_test_dnat.public_ip}",
	]
  }
  
  # Create CVM
  resource "tencentcloud_instance" "foo" {
	availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
	image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
	vpc_id           = "${data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id}"
	subnet_id         = "${data.tencentcloud_vpc_subnets.foo.instance_list.0.subnet_id}"
	system_disk_type  = "CLOUD_SSD"
  }
  
  # Add DNAT Entry
  resource "tencentcloud_dnat" "dev_dnat" {
	vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
	nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
	protocol     = "TCP"
	elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
	elastic_port = "80"
	private_ip   = "${tencentcloud_instance.foo.private_ip}"
	private_port = "9001"
	description	 = "test"
  }

data "tencentcloud_dnats" "multi_dnats" {
  nat_id           = "${tencentcloud_dnat.dev_dnat.nat_id}"
  vpc_id         = "${tencentcloud_dnat.dev_dnat.vpc_id}"
}
`
