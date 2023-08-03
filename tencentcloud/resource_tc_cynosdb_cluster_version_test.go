package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixCynosdbClusterVersionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterVersion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_version.cluster_version", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_version.cluster_version",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterVersion = `

resource "tencentcloud_cynosdb_cluster_version" "cluster_version" {
	cluster_id    = "cynosdbmysql-bws8h88b"
	cynos_version = "2.1.10"
}

`
