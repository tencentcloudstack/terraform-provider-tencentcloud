package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmClusterAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmClusterAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcm_cluster_attachment.cluster_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcm_cluster_attachment.cluster_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcmClusterAttachment = `

resource "tencentcloud_tcm_cluster_attachment" "cluster_attachment" {
  mesh_id = "mesh-xxxxxxxx"
  cluster_list {
		cluster_id = "cls-xxxxxxxx"
		region = "ap-shanghai"
		role = "REMOTE"
		vpc_id = "vpc-xxxxxxxx"
		subnet_id = "subnet-xxxxxxx"
		type = "TKE or EKS"

  }
}

`
