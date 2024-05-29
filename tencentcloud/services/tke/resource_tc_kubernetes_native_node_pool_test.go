package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesNativeNodePoolResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTkeNativeNodePool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "name", "tf-native-node-pool1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "type", "Native"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.0.name", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.0.value", "test2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.0.key", "product"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.0.value", "coderider"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.resource_type", "machine"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.0.key", "keep-test-np1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.0.value", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.1.key", "keep-test-np2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.1.value", "test2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "unschedulable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.system_disk.0.disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.system_disk.0.disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.instance_types.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.auto_repair", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.instance_charge_prepaid.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.instance_charge_prepaid.0.period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.instance_charge_prepaid.0.renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.0.nameservers.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.0.hosts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.0.kernel_args.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.host_name_pattern", "aaa{R:3}"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.kubelet_args.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.0.pre_init", "ZWNobyBoZWxsbw=="),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.0.post_init", "ZWNobyBoZWxsbw=="),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.runtime_root_dir", "/var/lib/docker"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.enable_autoscaling", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.replicas", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.internet_accessible.0.max_bandwidth_out", "50"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.internet_accessible.0.charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.file_system", "ext4"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.mount_target", "/var/lib/containerd"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.auto_format_and_mount", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.key_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "annotations.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "annotations.0.name"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "annotations.0.value"),
				),
			},
			{
				Config: testAccTkeNativeNodePoolUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "id"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "name", "tf-native-node-pool2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.0.name", "test11"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "labels.0.value", "test21"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "taints.0.effect", "NoExecute"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.1.key", "keep-test-np3"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "tags.0.tags.1.value", "test3"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.scaling.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.scaling.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.scaling.0.max_replicas", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.scaling.0.create_policy", "ZoneEquality"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.security_group_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.auto_repair", "false"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.0.hosts.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.management.0.kernel_args.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.kubelet_args.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.0.pre_init", "ZWNobyBoZWxsb3dvcmxk"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.lifecycle.0.post_init", "ZWNobyBoZWxsb3dvcmxk"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.enable_autoscaling", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.replicas"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.disk_size", "60"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.data_disks.0.auto_format_and_mount", "true"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "native.0.key_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_kubernetes_native_node_pool.native_node_pool_test", "annotations.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_kubernetes_native_node_pool.native_node_pool_test",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccNativeNodePoolDependentResource = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "vpc_cidr" {
  default = "172.16.0.0/16"
}

variable "subnet_cidr1" {
  default = "172.16.0.0/20"
}

variable "subnet_cidr2" {
  default = "172.16.16.0/20"
}

variable "tke_cidr_a" {
  default = [
    "10.31.0.0/23",
    "10.31.2.0/24",
    "10.31.3.0/24",
    "10.31.16.0/24",
    "10.31.32.0/24"
  ]
}

variable "default_img_id" {
  default = "img-2lr9q49h"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "tf-tke-vpc"
  cidr_block = var.vpc_cidr
}

resource "tencentcloud_subnet" "subnet1" {
  name              = "tf_tke_vpc_subnet1"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr1
  is_multicast      = false
}

resource "tencentcloud_subnet" "subnet2" {
  name              = "tf_tke_vpc_subnet2"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  cidr_block        = var.subnet_cidr2
  is_multicast      = false
}

data "tencentcloud_security_groups" "security_groups1" {
  name = "keep-tke"
}

data "tencentcloud_security_groups" "security_groups2" {
  name = "keep-reject-all"
}

locals {
  vpc_id     = tencentcloud_vpc.vpc.id
  subnet_id1 = tencentcloud_subnet.subnet1.id
  subnet_id2 = tencentcloud_subnet.subnet2.id

  sg_id1 = data.tencentcloud_security_groups.security_groups1.security_groups.0.security_group_id
  sg_id2 = data.tencentcloud_security_groups.security_groups2.security_groups.0.security_group_id
}

resource "tencentcloud_kubernetes_cluster" "kubernetes_cluster" {
  vpc_id                          = local.vpc_id
  cluster_cidr                    = var.tke_cidr_a.0
  cluster_max_pod_num             = 32
  cluster_name                    = "tf-tke-cluster-test"
  cluster_desc                    = "test cluster desc"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_domain         = "tf.cluster-internet.com"
  cluster_intranet                = false
  cluster_intranet_domain         = "tf.cluster-intranet.com"
  cluster_version                 = "1.22.5"
  cluster_os                      = "tlinux2.2(tkernel3)x86_64"
  cluster_level                   = "L5"
  auto_upgrade_cluster_level      = true
  cluster_internet_security_group = local.sg_id1
  node_name_type                  = "hostname"
  cluster_deploy_type             = "MANAGED_CLUSTER"
  unschedulable                   = 0
  cluster_subnet_id               = local.subnet_id1
}

resource "tencentcloud_key_pair" "key_pair1" {
  key_name = "tke_native_np_key1"
}

resource "tencentcloud_key_pair" "key_pair2" {
  key_name = "tke_native_np_key2"
}

locals {
  ssh1 = tencentcloud_key_pair.key_pair1.id
  ssh2 = tencentcloud_key_pair.key_pair2.id
}
`

const testAccTkeNativeNodePool = testAccNativeNodePoolDependentResource + `

resource "tencentcloud_kubernetes_native_node_pool" "native_node_pool_test" {
  cluster_id = tencentcloud_kubernetes_cluster.kubernetes_cluster.id
  name       = "tf-native-node-pool1"
  type       = "Native"

  labels {
    name  = "test1"
    value = "test2"
  }

  taints {
    key    = "product"
    value  = "coderider"
    effect = "NoSchedule"
  }

  tags {
    resource_type = "machine"
    tags {
      key   = "keep-test-np1"
      value = "test1"
    }
    tags {
      key   = "keep-test-np2"
      value = "test2"
    }
  }

  deletion_protection = false
  unschedulable       = false

  native {
    subnet_ids           = [local.subnet_id1]
    instance_charge_type = "PREPAID"
    system_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }
    instance_types     = ["SA2.MEDIUM2"]
    security_group_ids = [local.sg_id1]
    auto_repair        = false
    instance_charge_prepaid {
      period     = 1
      renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    }
    management {
      nameservers = ["183.60.83.19", "183.60.82.98"]
      hosts       = ["192.168.2.42 static.fake.com"]
      kernel_args = ["kernel.pid_max=65535"]
    }
    host_name_pattern = "aaa{R:3}"
    kubelet_args      = ["allowed-unsafe-sysctls=net.core.somaxconn"]
    lifecycle {
      pre_init  = "ZWNobyBoZWxsbw=="
      post_init = "ZWNobyBoZWxsbw=="
    }
    runtime_root_dir   = "/var/lib/docker"
    enable_autoscaling = false
    replicas           = 1
    internet_accessible {
      max_bandwidth_out = 50
      charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    }
    data_disks {
      disk_type             = "CLOUD_PREMIUM"
      file_system           = "ext4"
      disk_size             = 50
      mount_target          = "/var/lib/containerd"
      auto_format_and_mount = false
    }
    key_ids = [tencentcloud_key_pair.key_pair1.id]
  }

  annotations {
    name  = "cluster-autoscaler.kubernetes.io/scale-down-disabled"
    value = "true"
  }
}
`

const testAccTkeNativeNodePoolUpdate = testAccNativeNodePoolDependentResource + `

resource "tencentcloud_kubernetes_native_node_pool" "native_node_pool_test" {
  cluster_id = tencentcloud_kubernetes_cluster.kubernetes_cluster.id
  name       = "tf-native-node-pool2"
  type       = "Native"

  labels {
    name  = "test11"
    value = "test21"
  }

  taints {
    key    = "product"
    value  = "coderider"
    effect = "NoExecute"
  }

  tags {
    resource_type = "machine"
    tags {
      key   = "keep-test-np1"
      value = "test1"
    }
    tags {
      key   = "keep-test-np3"
      value = "test3"
    }
  }

  deletion_protection = false
  unschedulable       = false

  native {
    scaling {
      min_replicas  = 1
      max_replicas  = 2
      create_policy = "ZoneEquality"
    }
    subnet_ids           = [local.subnet_id1, local.subnet_id2]
    instance_charge_type = "PREPAID"
    system_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }
    instance_types     = ["SA2.MEDIUM2"]
    security_group_ids = [local.sg_id1, local.sg_id2]
    auto_repair        = false
    instance_charge_prepaid {
      period     = 1
      renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    }
    management {
      nameservers = ["183.60.83.19", "183.60.82.98"]
      hosts       = ["192.168.2.42 static.fake.com", "192.168.2.42 static.fake.com2"]
      kernel_args = ["kernel.pid_max=65535", "fs.file-max=400000"]
    }
    host_name_pattern = "aaa{R:3}"
    kubelet_args      = ["allowed-unsafe-sysctls=net.core.somaxconn", "root-dir=/var/lib/test"]
    lifecycle {
      pre_init  = "ZWNobyBoZWxsb3dvcmxk"
      post_init = "ZWNobyBoZWxsb3dvcmxk"
    }
    runtime_root_dir   = "/var/lib/docker"
    enable_autoscaling = true
    replicas           = 2
    internet_accessible {
      max_bandwidth_out = 50
      charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    }
    data_disks {
      disk_type             = "CLOUD_PREMIUM"
      file_system           = "ext4"
      disk_size             = 60
      mount_target          = "/var/lib/containerd"
      auto_format_and_mount = true
    }
    key_ids = [local.ssh1, local.ssh2]
  }

  annotations {
    name  = "cluster-autoscaler.kubernetes.io/scale-down-disabled"
    value = "true"
  }
  annotations {
    name  = "node.tke.cloud.tencent.com/security-agent"
    value = "false"
  }
}
`
