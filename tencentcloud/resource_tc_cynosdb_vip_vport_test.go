package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbVipVportResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbVipVport,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_vip_vport.vip_vport", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_vip_vport.vip_vport",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbVipVport = `

resource "tencentcloud_cynosdb_vip_vport" "vip_vport" {
  cluster_id = "xxx"
  instance_grp_id = "xxx"
  vip = "xx.xx.xx.xx"
  vport = 5432
  db_type = "MYSQL"
  old_ip_reserve_hours = 0
}

`
