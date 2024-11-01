package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixKubernetesClusterMasterAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesClusterMasterAttachment,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_master_attachment.example", "id"),
			),
		}},
	})
}

const testAccKubernetesClusterMasterAttachment = `
resource "tencentcloud_kubernetes_cluster_master_attachment" "example" {
  cluster_id                  = "cls-fp5o961e"
  instance_id                 = "ins-7d6tpbyg"
  node_role                   = "MASTER_ETCD"
  enhanced_security_service   = true
  enhanced_monitor_service    = true
  enhanced_automation_service = true
  password                    = "Password@123"
  security_group_ids          = ["sg-hjs685q9"]

  master_config {
    mount_target      = "/var/data"
    docker_graph_path = "/var/lib/containerd"
    unschedulable     = 0
    labels {
      name  = "key"
      value = "value"
    }

    data_disk {
      file_system           = "ext4"
      auto_format_and_mount = true
      mount_target          = "/var/data"
      disk_partition        = "/dev/vdb"
    }

    extra_args {
      kubelet = ["root-dir=/root"]
    }

    taints {
      key    = "key"
      value  = "value"
      effect = "NoSchedule"
    }
  }
}
`
