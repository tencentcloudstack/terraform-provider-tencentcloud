package tke_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudKubernetesClusterMasterAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccKubernetesClusterMasterAttachment,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_kubernetes_cluster_master_attachment.kubernetes_cluster_master_attachment", "id")),
		}, {
			ResourceName:      "tencentcloud_kubernetes_cluster_master_attachment.kubernetes_cluster_master_attachment",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccKubernetesClusterMasterAttachment = `

resource "tencentcloud_kubernetes_cluster_master_attachment" "kubernetes_cluster_master_attachment" {
  extra_args = {
  }
  master_config = {
    labels = {
    }
    data_disks = {
    }
    extra_args = {
    }
    gpu_args = {
    }
    taints = {
    }
  }
}
`
