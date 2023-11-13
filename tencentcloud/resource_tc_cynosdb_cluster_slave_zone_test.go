package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbCluster_slave_zoneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbCluster_slave_zone,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbCluster_slave_zone = `

resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = "cynosdbmysql-xxxxxxxx"
  slave_zone = "ap-guangzhou-3"
}

`
