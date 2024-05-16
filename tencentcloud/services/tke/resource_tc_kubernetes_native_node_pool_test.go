package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNativeNodePoolResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNativeNodePool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_native_node_pool.native_node_pool", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_native_node_pool.native_node_pool", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "name", "tf-native-node-pool1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "type", "Native"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.0.name", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.0.value", "test2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.0.key", "product"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.0.value", "coderider"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.resource_type", "machine"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.0.key", "keep-test-np1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.0.value", "test1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.1.key", "keep-test-np2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.1.value", "test2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "deletion_protection", "true"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "unschedulable", "false"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.instance_charge_type", "PREPAID"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.system_disk.0.disk_type", "CLOUD_SSD"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.system_disk.0.disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.instance_types.0。#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.security_group_ids.0.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.auto_repair", "false"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.instance_charge_prepaid.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.instance_charge_prepaid.0.period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.instance_charge_prepaid.0.renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.0.nameservers.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.0.hosts.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.0.kernel_args.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.host_name_pattern", "aaa{R:3}"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.kubelet_args.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.0.pre_init", "ZWNobyBoZWxsbw=="),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.0.post_init", "ZWNobyBoZWxsbw=="),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.runtime_root_dir", "/var/lib/docker"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.enable_autoscaling", "false"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.replicas", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.internet_accessible.0.max_bandwidth_out", "50"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.internet_accessible.0.charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.disk_type", "CLOUD_PREMIUM"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.file_system", "ext4"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.disk_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.mount_target", "/var/lib/containerd"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.auto_format_and_mount", "false"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.key_ids.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "annotations.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_native_node_pool.native_node_pool", "annotations.0.name"),
					resource.TestCheckResourceAttrSet("tencentcloud_native_node_pool.native_node_pool", "annotations.0.value"),
				),
			},
			{
				Config: testAccNativeNodePoolUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_native_node_pool.native_node_pool", "id"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "name", "tf-native-node-pool2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.0.name", "test11"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "labels.0.value", "test21"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "taints.0.effect", "NoExecute"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.1.key", "keep-test-np3"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "tags.0.tags.1.value", "test3"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.scaling.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.scaling.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.scaling.0.max_replicas", "10"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.scaling.0.create_policy", "ZoneEquality"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.security_group_ids.0.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.auto_repair", "true"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.0.hosts.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.management.0.kernel_args.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.health_check_policy_name", "policy1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.kubelet_args.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.0.pre_init", "ZWNobyBoZWxsb3dvcmxk"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.lifecycle.0.post_init", "ZWNobyBoZWxsb3dvcmxk"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.enable_autoscaling", "true"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.replicas", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.disk_size", "60"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.data_disks.0.auto_format_and_mount", "true"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "native.0.key_ids.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_native_node_pool.native_node_pool", "annotations.#", "2"),
				),
			},
			{
				ResourceName:      "tencentcloud_native_node_pool.native_node_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccNativeNodePool = `

resource "tencentcloud_native_node_pool" "tf_native_node_pool" {
  cluster_id = "cls-g52vxxss"
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

  deletion_protection = true
  unschedulable       = false

  native {
    subnet_ids           = ["subnet-83dxrme2"]
    instance_charge_type = "PREPAID"
    system_disk {
      disk_type = "CLOUD_SSD"
      disk_size = 50
    }
    instance_types     = ["SA2.MEDIUM2"]
    security_group_ids = ["sg-3o46faav"]
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
      disk_type              = "CLOUD_PREMIUM"
      file_system            = "ext4"
      disk_size              = 50
      mount_target           = "/var/lib/containerd"
      auto_format_and_mount  = false
    }
    key_ids = ["skey-4woe2xtj"]
  }

  annotations {
    name  = "cluster-autoscaler.kubernetes.io/scale-down-disabled"
    value = "true"
  }
}
`

const testAccNativeNodePoolUpdate = `

resource "tencentcloud_native_node_pool" "tf_native_node_pool" {
  cluster_id = "cls-g52vxxss"
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
      max_replicas  = 10
      create_policy = "ZoneEquality"
    }
    subnet_ids           = ["subnet-83dxrme2", "subnet-qp8hcbdk"]
    instance_charge_type = "PREPAID"
    system_disk {
      disk_type             = "CLOUD_SSD"
      disk_size             = 50
    }
    instance_types     = ["SA2.MEDIUM2"]
    security_group_ids = ["sg-3o46faav", "sg-cm7fbbf3"]
    auto_repair        = true
    instance_charge_prepaid {
      period     = 1
      renew_flag = "NOTIFY_AND_MANUAL_RENEW"
    }
    management {
      nameservers = ["183.60.83.19", "183.60.82.98"]
      hosts       = ["192.168.2.42 static.fake.com", "192.168.2.42 static.fake.com2"]
      kernel_args = ["kernel.pid_max=65535", "fs.file-max=400000"]
    }
    health_check_policy_name = "policy1"
    host_name_pattern        = "aaa{R:3}"
    kubelet_args             = ["allowed-unsafe-sysctls=net.core.somaxconn", "root-dir=/var/lib/test"]
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
      disk_type              = "CLOUD_PREMIUM"
      file_system            = "ext4"
      disk_size              = 60
      mount_target           = "/var/lib/containerd"
      auto_format_and_mount  = true
    }
    key_ids = ["skey-4woe2xtj", "skey-qp83bhp1"]
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
