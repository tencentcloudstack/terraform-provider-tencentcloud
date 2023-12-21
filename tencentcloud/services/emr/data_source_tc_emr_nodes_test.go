package emr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudEMRNodes(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEMRNodes(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_emr_nodes.my_emr_nodes"),
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
