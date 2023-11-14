package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigDatabaseCDCResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseCDC,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_c_d_c.config_database_c_d_c", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_c_d_c.config_database_c_d_c",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseCDC = `

resource "tencentcloud_sqlserver_config_database_c_d_c" "config_database_c_d_c" {
  d_b_names = 
  modify_type = "enable"
  instance_id = "mssql-i1z41iwd"
}

`
