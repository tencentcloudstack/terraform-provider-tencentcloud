package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverGeneralCommunicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCommunication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_communication.general_communication", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_communication.general_communication",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverGeneralCommunication = `

resource "tencentcloud_sqlserver_general_communication" "general_communication" {
  instance_id = "Instance ID in the format of mssql-j8kv137v"
  rename_restore {
		old_name = "old_db_name"
		new_name = "new_db_name"

  }
}

`
