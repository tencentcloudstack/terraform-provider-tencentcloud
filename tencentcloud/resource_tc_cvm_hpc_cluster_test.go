package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudNeedFixCvmHpcClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmHpcCluster,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_hpc_cluster.hpc_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmHpcCluster = `

resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone = "ap-beijing-6"
  name = "terraform-test"
  remark = "create for test"
}

`
