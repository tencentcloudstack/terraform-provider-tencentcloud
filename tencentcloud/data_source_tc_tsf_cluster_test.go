package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfClusterDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_cluster.cluster")),
			},
		},
	})
}

const testAccTsfClusterDataSource = `

data "tencentcloud_tsf_cluster" "cluster" {
  cluster_id_list = 
  cluster_type = "C"
  search_word = ""
  disable_program_auth_check = false
  }

`
