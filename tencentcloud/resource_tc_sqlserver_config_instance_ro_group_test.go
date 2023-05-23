package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverConfigInstanceRoGroupResource_basic -v
func TestAccTencentCloudSqlserverConfigInstanceRoGroupResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigInstanceRoGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_instance_ro_group.config_instance_ro_group", "id"),
				),
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
  instance_id = "mssql-ds1xhnt9"
  read_only_group_id = "mssqlrg-cbya44fb"
  read_only_group_name = "keep-ro-group-customize"
  is_offline_delay = 1
  read_only_max_delay_time = 10
  min_read_only_in_group = 1
  weight_pairs {
	read_only_instance_id = "mssqlro-o6dv2ugx"
	read_only_weight = 50
  }
  auto_weight = 0
}
`
