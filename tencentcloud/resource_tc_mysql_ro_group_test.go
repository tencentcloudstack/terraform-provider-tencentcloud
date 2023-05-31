package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlRoGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_ro_group.ro_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlRoGroup = `

resource "tencentcloud_mysql_ro_group" "ro_group" {
  ro_group_id = ""
  ro_group_info {
		ro_group_name = ""
		ro_max_delay_time = 
		ro_offline_delay = 
		min_ro_in_group = 
		weight_mode = ""
		replication_delay_time = 

  }
  ro_weight_values {
		instance_id = ""
		weight = 

  }
  is_balance_ro_load = 
  replication_delay_time = 
}

`
