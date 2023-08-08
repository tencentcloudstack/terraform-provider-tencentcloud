package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testTkeClusterNodePoolName = "tencentcloud_kubernetes_node_pool"
var testTkeClusterNodePoolResourceKey = testTkeClusterNodePoolName + ".np_test"

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_node_pool
	resource.AddTestSweepers("tencentcloud_node_pool", &resource.Sweeper{
		Name: "tencentcloud_node_pool",
		F:    testNodePoolSweep,
	})
}

var nodePoolNameReg = regexp.MustCompile("^(mynodepool|np)")

func testNodePoolSweep(region string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cli, err := sharedClientForRegion(region)
	if err != nil {
		return err
	}
	client := cli.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	request := tke.NewDescribeClusterNodePoolsRequest()
	request.ClusterId = helper.String(defaultTkeClusterId)
	response, err := client.UseTkeClient().DescribeClusterNodePools(request)
	if err != nil {
		log.Printf("Query %s node pool fail: %s", defaultTkeClusterId, err.Error())
	}
	nodePools := response.Response.NodePoolSet
	if len(nodePools) == 0 {
		return nil
	}
	for i := range nodePools {
		poolId := *nodePools[i].NodePoolId
		poolName := nodePools[i].Name
		if poolName == nil {
			continue
		}

		if !nodePoolNameReg.MatchString(*poolName) {
			continue
		}
		err := service.DeleteClusterNodePool(ctx, defaultTkeClusterId, poolId, false)
		if err != nil {
			continue
		}
	}
	return nil
}

func TestAccTencentCloudKubernetesNodePoolResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeNodePoolCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.system_disk_size", "50"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.data_disk.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.internet_max_bandwidth_out", "10"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.cam_role_name", "TCB_QcsRole"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "taints.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "max_size", "6"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "min_size", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "desired_capacity", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "name", "mynodepool"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "unschedulable", "0"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "scaling_group_name", "asg_np_test"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "default_cooldown", "400"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "termination_policies.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "termination_policies.0", "OLDEST_INSTANCE"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_count", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "autoscaling_added_total", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "manually_added_total", "0"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "tags.keep-test-np1", "test1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "tags.keep-test-np2", "test2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.host_name", "12.123.0.0"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.host_name_style", "ORIGINAL"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.enhanced_security_service", "false"),
				),
			},
			{
				Config: testAccTkeNodePoolClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.system_disk_size", "100"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.data_disk.#", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.data_disk.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.data_disk.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.internet_max_bandwidth_out", "20"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.instance_charge_type", "SPOTPAID"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.spot_instance_type", "one-time"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.spot_max_price", "1000"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.cam_role_name", "TCB_QcsRole"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "max_size", "5"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "min_size", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "labels.test3", "test3"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "desired_capacity", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "name", "mynodepoolupdate"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "node_os", defaultTkeOSImageName),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "unschedulable", "0"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "scaling_group_name", "asg_np_test_changed"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "default_cooldown", "350"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "termination_policies.#", "1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "termination_policies.0", "NEWEST_INSTANCE"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "tags.keep-test-np1", "testI"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "tags.keep-test-np3", "testIII"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.security_group_ids.#", "2"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.host_name", "12.123.1.1"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.host_name_style", "UNIQUE"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.enhanced_security_service", "true"),
				),
			},
		},
	})
}

func TestAccTencentCloudKubernetesNodePoolResource_DiskEncrypt(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeNodePoolClusterEncrypt,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttr(testTkeClusterNodePoolResourceKey, "auto_scaling_config.0.data_disk.0.encrypt", "true"),
				),
			},
		},
	})
}

func TestAccTencentCloudKubernetesNodePoolResource_GPUInstance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTkeNodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeNodePoolClusterGpu,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTkeNodePoolExists,
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "cluster_id"),
					resource.TestCheckResourceAttrSet(testTkeClusterNodePoolResourceKey, "node_config.0.gpu_args.#"),
				),
			},
		},
	})
}

func testAccCheckTkeNodePoolDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterNodePoolResourceKey]
	if !ok {
		return fmt.Errorf("tke node pool %s is not found", testTkeClusterNodePoolResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke  node pool id is not set")
	}
	items := strings.Split(rs.Primary.ID, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id %s is broken", rs.Primary.ID)
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		if err.(*sdkErrors.TencentCloudSDKError).Code == "InternalError.UnexpectedInternal" {
			return nil
		}
		return err
	}
	if !has {
		return nil
	} else {
		return fmt.Errorf("tke node pool %s still exist", nodePoolId)
	}

}

func testAccCheckTkeNodePoolExists(s *terraform.State) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	rs, ok := s.RootModule().Resources[testTkeClusterNodePoolResourceKey]
	if !ok {
		return fmt.Errorf("tke node pool %s is not found", testTkeClusterNodePoolResourceKey)
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("tke node pool id is not set")
	}

	items := strings.Split(rs.Primary.ID, FILED_SP)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  %s is broken", rs.Primary.ID)
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		return err
	}
	if has {
		return nil
	} else {
		return fmt.Errorf("tke node pool %s query fail.", nodePoolId)
	}

}

const testAccTkeNodePoolClusterBasic = defaultProjectVariable + defaultImages + TkeDataSource + TkeDefaultNodeInstanceVar + `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

data "tencentcloud_security_groups" "sg" {
  name = "default"
}

data "tencentcloud_security_groups" "sg_as" {
  name = "keep-for-as"
}
`

const testAccTkeNodePoolCluster string = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "mynodepool"
  cluster_id = local.cluster_id
  max_size = 6
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 1
  enable_auto_scale    = true
  scaling_group_name	   = "asg_np_test"
  default_cooldown		   = 400
  termination_policies	   = ["OLDEST_INSTANCE"]
  scaling_group_project_id = var.default_project
  delete_keep_instance = false
  node_os="tlinux2.2(tkernel3)x86_64"

  auto_scaling_config {
    instance_type      = var.ins_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = [data.tencentcloud_security_groups.sg.security_groups[0].security_group_id]
    cam_role_name = "TCB_QcsRole"
    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false
	host_name                  = "12.123.0.0"
	host_name_style            = "ORIGINAL"
  }
  unschedulable = 0
  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  tags = {
    keep-test-np1 = "test1"
    keep-test-np2 = "test2"
  }

  node_config {
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}
`

const testAccTkeNodePoolClusterUpdate string = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "mynodepoolupdate"
  cluster_id = local.cluster_id
  max_size = 5
  min_size = 2
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 2
  enable_auto_scale    = false
  node_os = var.default_img
  scaling_group_project_id = var.default_project
  delete_keep_instance = false
  scaling_group_name 	   = "asg_np_test_changed"
  default_cooldown 		   = 350
  termination_policies 	   = ["NEWEST_INSTANCE"]
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.ins_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "100"
    security_group_ids = [data.tencentcloud_security_groups.sg.security_groups[0].security_group_id, data.tencentcloud_security_groups.sg_as.security_groups[0].security_group_id]
	instance_charge_type = "SPOTPAID"
    spot_instance_type = "one-time"
    spot_max_price = "1000"
    cam_role_name = "TCB_QcsRole"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      delete_with_instance = true
    }
    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 100
      delete_with_instance = true
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 20
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = true
    enhanced_monitor_service   = false
	host_name                  = "12.123.1.1"
	host_name_style            = "UNIQUE"

  }
  unschedulable = 0
  labels = {
    "test3" = "test3",
    "test2" = "test2",
  }
  
  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  tags = {
    keep-test-np1 = "testI"
    keep-test-np3 = "testIII"
  }

  node_config {
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
  }
}
`

const testAccTkeNodePoolClusterEncrypt = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "np_with_disk_encrypt"
  cluster_id = local.cluster_id
  max_size = 3
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 1
  enable_auto_scale    = true
  scaling_group_name	   = "encrypt_asg"
  default_cooldown		   = 400
  termination_policies	   = ["OLDEST_INSTANCE"]
  scaling_group_project_id = var.default_project
  delete_keep_instance = false
  node_os="tlinux2.2(tkernel3)x86_64"

  auto_scaling_config {
    instance_type      = var.ins_type
    cam_role_name      = "TCB_QcsRole"
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = [data.tencentcloud_security_groups.sg.security_groups[0].security_group_id]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      encrypt   = true
    }
    public_ip_assigned         = false
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

  }
  unschedulable = 0
}
`

const testAccTkeNodePoolClusterGpu string = testAccTkeNodePoolClusterBasic + `
resource "tencentcloud_kubernetes_node_pool" "np_test" {
  name = "gpu_args_node_pool"
  cluster_id = local.cluster_id
  max_size = 1
  min_size = 0
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 1
  enable_auto_scale    = false
  node_os = "tlinux3.1x86_64"
  scaling_group_project_id = var.default_project
  delete_keep_instance = false
  scaling_group_name 	   = "asg_np_test_changed_gpu"
  default_cooldown 		   = 350
  termination_policies 	   = ["NEWEST_INSTANCE"]
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = "GN6S.LARGE20"
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "100"
    security_group_ids = [data.tencentcloud_security_groups.sg.security_groups[0].security_group_id, data.tencentcloud_security_groups.sg_as.security_groups[0].security_group_id]
	instance_charge_type = "SPOTPAID"
    spot_instance_type = "one-time"
    spot_max_price = "1000"
    cam_role_name = "TCB_QcsRole"

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
      delete_with_instance = true
    }
    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 100
      delete_with_instance = true
    }

    public_ip_assigned         = false
    password                   = "test123#"
    enhanced_security_service  = true
    enhanced_monitor_service   = false
	host_name                  = "12.123.1.1"
	host_name_style            = "UNIQUE"

  }
  unschedulable = 0
  labels = {
    "test3" = "test3",
    "test2" = "test2",
  }
  
  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  tags = {
    keep-test-np1 = "testI"
    keep-test-np3 = "testIII"
  }

  node_config {
    extra_args = [
      "root-dir=/var/lib/kubelet"
    ]
	gpu_args {
      mig_enable = false
      driver = {
        name = "NVIDIA-Linux-x86_64-470.182.03.run"
        version = "470.182.03"
      }
      cuda = {
        name = "cuda_11.4.3_470.82.01_linux.run"
        version = "11.4.3"
      }
      cudnn = {
        name = "cudnn-11.4-linux-x64-v8.2.4.15.tgz"
        version = "8.2.4"
      }
    }
  }
}
`
