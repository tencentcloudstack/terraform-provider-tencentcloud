package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcmMeshDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcmMeshDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcm_mesh.mesh")),
			},
		},
	})
}

const testAccTcmMeshDataSource = `

data "tencentcloud_tcm_mesh" "mesh" {
  mesh_id = 
  mesh_name = 
  tags = 
  mesh_cluster = 
  }

`
