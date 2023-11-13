package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverConfigDatabaseMdfResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverConfigDatabaseMdf,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_config_database_mdf.config_database_mdf", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_config_database_mdf.config_database_mdf",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverConfigDatabaseMdf = `

resource "tencentcloud_sqlserver_config_database_mdf" "config_database_mdf" {
  d_b_names = 
  instance_id = "mssql-i1z41iwd"
}

`
