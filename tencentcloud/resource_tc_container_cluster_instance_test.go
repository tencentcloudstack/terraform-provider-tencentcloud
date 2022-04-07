package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// @Deprecated It has been deprecated and replaced by tencentcloud_kubernetes_scale_worker.
func testAccTencentCloudContainerClusterInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterInstanceDestroy,
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

func testAccCheckContainerClusterInstanceDestroy(s *terraform.State) error {
	return nil
}

const testAccTencentCloudContainerClusterInstanceConfig_basic = `
variable "my_vpc" {
  default = "` + defaultVpcId + `"
}

variable "my_subnet" {
  default = "` + defaultSubnetId + `"
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
  os_name           = "ubuntu16.04.1 LTSx86_64"
  bandwidth         = 1
  bandwidth_type    = "PayByHour"
  require_wan_ip    = 1
  subnet_id         = var.my_subnet
  is_vpc_gateway    = 0
  storage_size      = 0
  root_size         = 50
  goods_num         = 2
  password          = "Admin12345678"
  vpc_id            = var.my_vpc
  cluster_cidr      = "10.0.32.0/19"
  cvm_type          = "PayByHour"
  cluster_desc      = "foofoofoo"
  period            = 1
  zone_id           = 100003
  instance_type     = data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type
  instance_name     = "terraform-container-acc-test-vm"
  cluster_version   = "1.14.3"
}

resource "tencentcloud_container_cluster_instance" "bar_instance" {
  cluster_id        = tencentcloud_container_cluster.foo.id
  bandwidth         = 1
  bandwidth_type    = "PayByHour"
  require_wan_ip    = 1
  subnet_id         = var.my_subnet
  is_vpc_gateway    = 0
  storage_size      = 10
  root_size         = 50
  root_type 		    = "CLOUD_SSD"
  password          = "Admin12345678"
  cvm_type          = "PayByHour"
  period            = 1
  zone_id           = 100003
  instance_type     = data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type
  instance_name     = "terraform-container-acc-test-vm-add"
}
`
