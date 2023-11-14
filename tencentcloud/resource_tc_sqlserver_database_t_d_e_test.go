package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverDatabaseTDEResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverDatabaseTDE,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_database_t_d_e.database_t_d_e", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_database_t_d_e.database_t_d_e",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverDatabaseTDE = `

resource "tencentcloud_sqlserver_database_t_d_e" "database_t_d_e" {
  instance_id = "mssql-i1z41iwd"
  d_b_t_d_e_encrypt {
		d_b_name = ""
		encryption = ""

  }
}

`
