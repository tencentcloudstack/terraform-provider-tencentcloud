package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterZoneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterZone,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_zone.cluster_zone", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_zone.cluster_zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterZone = `

resource "tencentcloud_cynosdb_cluster_zone" "cluster_zone" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  old_zone = "ap-guangzhou-2"
  new_zone = "ap-guangzhou-3"
  is_in_maintain_period = "false"
}

`
