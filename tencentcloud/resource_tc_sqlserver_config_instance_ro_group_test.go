package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigInstanceRoGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceRoGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_ro_group.config_instance_ro_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_instance_ro_group.config_instance_ro_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigInstanceRoGroup = `

resource "tencentcloud_sqlserver_config_instance_ro_group" "config_instance_ro_group" {
  instance_id = "mssql-i1z41iwd"
  read_only_group_id = ""
  read_only_group_name = ""
  is_offline_delay = 
  read_only_max_delay_time = 
  min_read_only_in_group = 
  weight_pairs {
		read_only_instance_id = ""
		read_only_weight = 

  }
  auto_weight = 
  balance_weight = 
}

`
