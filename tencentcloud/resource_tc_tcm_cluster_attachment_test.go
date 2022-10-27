package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTcmClusterAttachment_basic -v
func TestAccTencentCloudTcmClusterAttachment_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusterAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmClusterAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterAttachmentExists("tencentcloud_tcm_cluster_attachment.basic"),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "mesh_id", defaultMeshId),
					resource.TestCheckResourceAttr("tencentcloud_tcm_cluster_attachment.basic", "cluster_list.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcm_cluster_attachment.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClusterAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcm_cluster_attachment" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		meshId := ids[0]
		clusterId := ids[1]

		mesh, err := service.DescribeTcmMesh(ctx, meshId)
		if err != nil {
			return err
		}

		if mesh != nil {
			if len(mesh.Mesh.ClusterList) > 0 {
				for _, v := range mesh.Mesh.ClusterList {
					if *v.ClusterId == clusterId {
						return fmt.Errorf("clusterList %s still exists", rs.Primary.ID)
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckClusterAttachmentExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 2 {
			return fmt.Errorf("id is broken, id is %s", rs.Primary.ID)
		}
		meshId := ids[0]
		clusterId := ids[1]

		service := TcmService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		mesh, err := service.DescribeTcmMesh(ctx, meshId)
		if err != nil {
			return err
		}

		if mesh == nil {
			return fmt.Errorf("mesh %s is not found", rs.Primary.ID)
		}
		if len(mesh.Mesh.ClusterList) > 0 {
			for _, v := range mesh.Mesh.ClusterList {
				if *v.ClusterId == clusterId {
					return nil
				}
			}
			return fmt.Errorf("mesh clusterList %s is not found", rs.Primary.ID)
		} else {
			return fmt.Errorf("clusterList %s is not found", rs.Primary.ID)
		}
	}
}

const testAccTcmClusterAttachmentVar = `
variable "cluster_id" {
  default = "` + defaultMeshClusterId + `"
}
variable "mesh_id" {
  default = "` + defaultMeshId + `"
}
variable "vpc_id" {
  default = "` + defaultMeshVpcId + `"
}
variable "subnet_id" {
  default = "` + defaultMeshSubnetId + `"
}
`

const testAccTcmClusterAttachment = testAccTcmClusterAttachmentVar + `

resource "tencentcloud_tcm_cluster_attachment" "basic" {
  mesh_id = var.mesh_id
  cluster_list {
    cluster_id = var.cluster_id
    region = "ap-guangzhou"
    role = "REMOTE"
    vpc_id = var.vpc_id
    subnet_id = var.subnet_id
    type = "TKE"
  }
}

`
