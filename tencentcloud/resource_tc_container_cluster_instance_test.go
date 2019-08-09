package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudContainerClusterInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudContainerClusterInstanceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_container_cluster_instance.bar_instance"),
					checkContainerClusterInstancesAllNormal("tencentcloud_container_cluster.foo"),
				),
			},
		},
	})
}

const testAccTencentCloudContainerClusterInstanceConfig_basic = `
variable "my_vpc" {
   default = "`+DefaultVpcId+`"
}

variable "my_subnet" {
  default = "`+DefaultSubnetId+`"
}

resource "tencentcloud_container_cluster" "foo" {
 cluster_name = "terraform-acc-test-inst"
 cpu    = 1
 mem    = 1
 os_name   = "ubuntu16.04.1 LTSx86_64"
 bandwidth  = 1
 bandwidth_type = "PayByHour"
 require_wan_ip   = 1
 subnet_id  = "${var.my_subnet}"
 is_vpc_gateway = 0
 storage_size = 0
 root_size  = 50
 root_type = "CLOUD_SSD"
 goods_num  = 1
 password  = "Admin12345678#!"
 vpc_id   = "${var.my_vpc}"
 cluster_cidr = "10.0.0.0/19"
 cvm_type  = "PayByHour"
 cluster_desc = "foofoofoo"
 period   = 1
 zone_id   = 100003
 instance_type = "S4.SMALL1"
 mount_target = ""
 docker_graph_path = ""
 instance_name = "terraform-container-acc-test-vm"
 cluster_version = "1.7.8"
}

resource "tencentcloud_container_cluster_instance" "bar_instance" {
 cpu    = 1
 mem    = 1
 bandwidth  = 1
 bandwidth_type = "PayByHour"
 require_wan_ip   = 1
 is_vpc_gateway = 0
 storage_size = 10
 root_size  = 50
 root_type = "CLOUD_SSD"
 password  = "Admin12345678"
 cvm_type  = "PayByHour"
 period   = 1
 zone_id   = 100003
 instance_type = "CVM.S3"
 mount_target = "/data"
 docker_graph_path = ""
 subnet_id  = "${var.my_subnet}"
 cluster_id = "${tencentcloud_container_cluster.foo.id}"
}
`
