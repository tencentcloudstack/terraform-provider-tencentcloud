package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudDataSourceContainerClusterInstances(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_container_cluster_instances.foo_instance"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceContainerClusterInstancesConfig_basic = `
variable "my_vpc" {
   default = "`+DefaultVpcId+`"
}

variable "my_subnet" {
  default = "`+DefaultSubnetId+`"
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_container_cluster" "foo" {
  cluster_name      = "terraform-acc-test"
  cpu               = 1
  mem               = 1
  os_name           = "ubuntu16.04.1 LTSx86_64"
  bandwidth         = 1
  bandwidth_type    = "PayByHour"
  require_wan_ip    = 1
  subnet_id         = "${var.my_subnet}"
  is_vpc_gateway    = 0
  storage_size      = 0
  root_size         = 50
  root_type         = "CLOUD_SSD"
  goods_num         = 1
  password          = "Admin12345678"
  vpc_id            = "${var.my_vpc}"
  cluster_cidr      = "10.0.0.0/19"
  cvm_type          = "PayByHour"
  cluster_desc      = "foofoofoo"
  period            = 1
  zone_id           = 100003
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  mount_target      = ""
  docker_graph_path = ""
  instance_name     = "terraform-container-acc-test-vm"
  cluster_version   = "1.7.8"
}


data "tencentcloud_container_cluster_instances" "foo_instance" {
    cluster_id = "${tencentcloud_container_cluster.foo.id}"
}
`
