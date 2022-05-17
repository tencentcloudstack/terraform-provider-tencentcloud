package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudEMRNodes(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEMRNodes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_emr_nodes.my_emr_nodes"),
					resource.TestCheckResourceAttr("data.tencentcloud_emr_nodes.my_emr_nodes", "nodes.#", "1"),
				),
			},
		},
	})
}

func testAccEMRNodes() string {
	return testEmrBasic + `
data "tencentcloud_emr_nodes" "my_emr_nodes" {
  node_flag="master"
  instance_id=tencentcloud_emr_cluster.emrrrr.instance_id
}
`
}
