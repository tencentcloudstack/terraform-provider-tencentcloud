package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverGeneralBackupsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralBackups,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_backups.general_backups", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_backups.general_backups",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverGeneralBackups = `

resource "tencentcloud_sqlserver_general_backups" "general_backups" {
  strategy = 0
  d_b_names = 
  instance_id = "mssql-i1z41iwd"
  backup_name = "bk_name"
}

`
