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
resource "tencentcloud_vpc" "my_vpc" {
  cidr_block = "10.6.0.0/16"
  name       = "terraform_vpc_test"
}

resource "tencentcloud_subnet" "my_subnet" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  availability_zone = "ap-guangzhou-3"
  name              = "terraform_test_subnet"
  cidr_block        = "10.6.0.0/24"
}

resource "tencentcloud_container_cluster" "foo" {
 cluster_name = "bar"
 cpu    = 1
 mem    = 1
 os_name   = "ubuntu16.04.1 LTSx86_64"
 bandwidth  = 1
 bandwidth_type = "PayByHour"
 require_wan_ip   = 1
 subnet_id  = "${tencentcloud_subnet.my_subnet.id}"
 is_vpc_gateway = 0
 storage_size = 0
 root_size  = 50
 goods_num  = 1
 password  = "Admin12345678"
 vpc_id   = "${tencentcloud_vpc.my_vpc.id}"
 cluster_cidr = "10.0.0.0/19"
 cvm_type  = "PayByHour"
 cluster_desc = "foofoofoo"
 period   = 1
 zone_id   = 100003
 instance_type = "S2.SMALL1"
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
 password  = "Admin12345678"
 cvm_type  = "PayByHour"
 period   = 1
 zone_id   = 100003
 instance_type = "CVM.S2"
 mount_target = "/data"
 docker_graph_path = ""
 subnet_id  = "${tencentcloud_subnet.my_subnet.id}"
 cluster_id = "${tencentcloud_container_cluster.foo.id}"
}
`
