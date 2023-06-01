package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRoGroupResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.min_ro_in_group"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.replication_delay_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_group_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_max_delay_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.ro_offline_delay"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_group_info.0.weight_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.0.instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group.ro_group", "ro_weight_values.0.weight"),
				),
			},
		},
	})
}

const testAccMysqlRoGroup = `

resource "tencentcloud_mysql_ro_group" "ro_group" {
	instance_id = "cdb-e8i766hx"
	ro_group_id = "cdbrg-f49t0gnj"
	ro_group_info {
	  ro_group_name          = "keep-ro"
	  ro_max_delay_time      = 1
	  ro_offline_delay       = 1
	  min_ro_in_group        = 1
	  weight_mode            = "custom"
	  # replication_delay_time = 1
	}
	ro_weight_values {
	  instance_id = "cdbro-f49t0gnj"
	  weight      = 10
	}
	is_balance_ro_load = 1
}

`
