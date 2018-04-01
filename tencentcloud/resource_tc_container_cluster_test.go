package tencentcloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudContainerClusters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudContainerClustersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_container_cluster.foo"),
				),
			},
		},
	})
}

const testAccTencentCloudContainerClustersConfig_basic = `
resource "tencentcloud_vpc" "my_vpc" {
  cidr_block = "10.6.0.0/16"
  name       = "tf_vpc_test"
}

resource "tencentcloud_subnet" "my_subnet" {
  vpc_id = "${tencentcloud_vpc.my_vpc.id}"
  availability_zone = "ap-guangzhou-3"
  name              = "tf_test_subnet"
  cidr_block        = "10.6.0.0/24"
}

resource "tencentcloud_container_cluster" "foo" {
 cluster_name = "terraform-acc-test"
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
 instance_name = "bar-vm"
 cluster_version = "1.7.8"
}
`

