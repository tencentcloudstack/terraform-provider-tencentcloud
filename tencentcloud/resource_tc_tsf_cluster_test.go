package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfCluster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_cluster.cluster", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_cluster.cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfCluster = `

resource "tencentcloud_tsf_cluster" "cluster" {
  cluster_name = ""
  cluster_type = ""
  vpc_id = ""
  cluster_c_i_d_r = ""
  cluster_desc = ""
  tsf_region_id = ""
  tsf_zone_id = ""
  subnet_id = ""
  cluster_version = ""
  max_node_pod_num = 
  max_cluster_service_num = 
  program_id = ""
  kubernete_api_server = ""
  kubernete_native_type = ""
  kubernete_native_secret = ""
  program_id_list = 
  tags = {
    "createdBy" = "terraform"
  }
}

`
