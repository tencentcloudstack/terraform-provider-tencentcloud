package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigDatabaseCTResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseCT,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_c_t.config_database_c_t", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_c_t.config_database_c_t",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseCT = `

resource "tencentcloud_sqlserver_config_database_c_t" "config_database_c_t" {
  d_b_names = 
  modify_type = "enable"
  instance_id = "mssql-i1z41iwd"
  change_retention_day = 7
}

`
