package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfClusterDataSource_basic -v
func TestAccTencentCloudTsfClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfClusterDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_cluster.cluster"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.0.cluster_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.0.delete_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_cluster.cluster", "result.0.content.0.vpc_id"),
				),
			},
		},
	})
}

const testAccTsfClusterDataSourceVar = `
variable "cluster_id" {
	default = "` + defaultTsfClustId + `"
}
`

const testAccTsfClusterDataSource = testAccTsfClusterDataSourceVar + `

data "tencentcloud_tsf_cluster" "cluster" {
	cluster_id_list = [var.cluster_id]
	cluster_type = "V"
	# search_word = ""
	disable_program_auth_check = true
}

`
