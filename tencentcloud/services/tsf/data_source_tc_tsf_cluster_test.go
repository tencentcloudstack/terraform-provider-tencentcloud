package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfClusterDataSource_basic -v
func TestAccTencentCloudTsfClusterDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfClusterDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_cluster.cluster"),
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
	default = "` + tcacctest.DefaultTsfClustId + `"
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
