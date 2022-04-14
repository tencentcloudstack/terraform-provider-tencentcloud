package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// @Deprecated It has been deprecated and replaced by tencentcloud_kubernetes_cluster.
func testAccTencentCloudContainerCluster_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudContainerClusterConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_container_cluster.foo"),
					checkContainerClusterInstancesAllNormal("tencentcloud_container_cluster.foo"),
				),
			},
		},
	})
}

func testAccCheckContainerClusterDestroy(s *terraform.State) error {
	return nil
}

// For ordinary usage, it doesn't require all nodes in a cluster to be in normal state.
// But for acceptance test, it only has a single node and should be in normal state otherwise
// will cause resource leak such as vpc, subnet and vm resources, these leakage will block
// subsequential acceptance test, hence here we need to do such check to ensure cluster node
// is in an expected state.
func checkContainerClusterInstancesAllNormal(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Container cluster ID is not set")
		}
		ctx := context.TODO()
		service := TkeService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		var master []InstanceInfo
		var worker []InstanceInfo
		err := resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
			var e error
			master, worker, e = service.DescribeClusterInstances(ctx, rs.Primary.ID)
			if e != nil {
				return retryError(e)
			}
			allRunning := true
			for _, n := range master {
				if n.InstanceState != "running" {
					allRunning = false
					break
				}
			}
			for _, n := range worker {
				if n.InstanceState != "running" {
					allRunning = false
					break
				}
			}
			if allRunning {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("No all state is running"))
		})
		return err
	}
}

const testAccTencentCloudContainerClusterConfig_basic = `
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
  cluster_cidr      = "10.0.0.0/19"
  cvm_type          = "PayByHour"
  cluster_desc      = "foofoofoo"
  period            = 1
  zone_id           = 100003
  instance_type     = data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type
  mount_target      = ""
  docker_graph_path = ""
  instance_name     = "terraform-container-acc-test-vm"
  cluster_version   = "1.14.3"
}
`
